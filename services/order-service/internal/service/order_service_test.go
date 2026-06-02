package service_test

import (
	"errors"
	"testing"

	"github.com/448267450/food-delivery-platform/services/order-service/internal/model"
	"github.com/448267450/food-delivery-platform/services/order-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==================== Mock Repository ====================

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(id uint) (*model.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByRestaurantID(restaurantID uint) ([]model.Order, error) {
	args := m.Called(restaurantID)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateStatus(id uint, status model.OrderStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) Update(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

// ==================== Create Order Tests ====================

func TestCreateOrder_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	req := &service.CreateOrderRequest{
		UserID:       1,
		RestaurantID: 1,
		Address:      "456 Oak Ave, Austin TX",
		Note:         "No onions",
		Items: []service.OrderItemRequest{
			{MenuItemID: 1, Name: "Cheeseburger", Price: 12.99, Quantity: 2},
			{MenuItemID: 2, Name: "Coke", Price: 2.99, Quantity: 1},
		},
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Order")).Return(nil)

	order, err := svc.CreateOrder(req)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, model.StatusPending, order.Status)
	// Total = (12.99 * 2) + (2.99 * 1) = 28.97
	assert.Equal(t, 28.97, order.TotalPrice)
	assert.Len(t, order.Items, 2)
}

func TestCreateOrder_DatabaseError(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	req := &service.CreateOrderRequest{
		UserID:       1,
		RestaurantID: 1,
		Address:      "456 Oak Ave",
		Items: []service.OrderItemRequest{
			{MenuItemID: 1, Name: "Burger", Price: 10.00, Quantity: 1},
		},
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Order")).Return(errors.New("db connection lost"))

	order, err := svc.CreateOrder(req)

	assert.Error(t, err)
	assert.Nil(t, order)
}

// ==================== State Machine Tests ====================

// Test all valid transitions in sequence
func TestUpdateOrderStatus_ValidTransitions(t *testing.T) {
	validFlow := []struct {
		from model.OrderStatus
		to   model.OrderStatus
	}{
		{model.StatusPending, model.StatusPaid},
		{model.StatusPaid, model.StatusPreparing},
		{model.StatusPreparing, model.StatusOutForDelivery},
		{model.StatusOutForDelivery, model.StatusDelivered},
	}

	for _, tc := range validFlow {
		t.Run(string(tc.from)+"->"+string(tc.to), func(t *testing.T) {
			mockRepo := new(MockOrderRepository)
			svc := service.NewOrderService(mockRepo)

			order := &model.Order{ID: 1, Status: tc.from}
			mockRepo.On("FindByID", uint(1)).Return(order, nil)
			mockRepo.On("UpdateStatus", uint(1), tc.to).Return(nil)

			result, err := svc.UpdateOrderStatus(1, &service.UpdateStatusRequest{Status: tc.to})

			assert.NoError(t, err)
			assert.Equal(t, tc.to, result.Status)
		})
	}
}

// Test invalid transitions are rejected
func TestUpdateOrderStatus_InvalidTransitions(t *testing.T) {
	invalidTransitions := []struct {
		from model.OrderStatus
		to   model.OrderStatus
	}{
		{model.StatusPending, model.StatusDelivered},   // skip steps
		{model.StatusPending, model.StatusPreparing},   // skip PAID
		{model.StatusDelivered, model.StatusPending},   // terminal state
		{model.StatusCancelled, model.StatusPaid},      // terminal state
		{model.StatusPreparing, model.StatusCancelled}, // too late to cancel
	}

	for _, tc := range invalidTransitions {
		t.Run(string(tc.from)+"->"+string(tc.to), func(t *testing.T) {
			mockRepo := new(MockOrderRepository)
			svc := service.NewOrderService(mockRepo)

			order := &model.Order{ID: 1, Status: tc.from}
			mockRepo.On("FindByID", uint(1)).Return(order, nil)

			result, err := svc.UpdateOrderStatus(1, &service.UpdateStatusRequest{Status: tc.to})

			assert.Error(t, err)
			assert.Nil(t, result)
			// UpdateStatus should never reach the DB for invalid transitions
			mockRepo.AssertNotCalled(t, "UpdateStatus")
		})
	}
}

// ==================== Cancel Order Tests ====================

func TestCancelOrder_FromPending(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	order := &model.Order{ID: 1, Status: model.StatusPending}
	mockRepo.On("FindByID", uint(1)).Return(order, nil)
	mockRepo.On("UpdateStatus", uint(1), model.StatusCancelled).Return(nil)

	err := svc.CancelOrder(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCancelOrder_FromPaid(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	order := &model.Order{ID: 1, Status: model.StatusPaid}
	mockRepo.On("FindByID", uint(1)).Return(order, nil)
	mockRepo.On("UpdateStatus", uint(1), model.StatusCancelled).Return(nil)

	err := svc.CancelOrder(1)

	assert.NoError(t, err)
}

func TestCancelOrder_AlreadyPreparing(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	// Once restaurant starts preparing, cannot cancel
	order := &model.Order{ID: 1, Status: model.StatusPreparing}
	mockRepo.On("FindByID", uint(1)).Return(order, nil)

	err := svc.CancelOrder(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel order")
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

func TestCancelOrder_AlreadyDelivered(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	order := &model.Order{ID: 1, Status: model.StatusDelivered}
	mockRepo.On("FindByID", uint(1)).Return(order, nil)

	err := svc.CancelOrder(1)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

func TestCancelOrder_NotFound(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.CancelOrder(999)

	assert.Error(t, err)
	assert.Equal(t, "order not found", err.Error())
}

// ==================== Query Tests ====================

func TestGetOrdersByUser_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo)

	orders := []model.Order{
		{ID: 1, UserID: 1, Status: model.StatusDelivered},
		{ID: 2, UserID: 1, Status: model.StatusPending},
	}

	mockRepo.On("FindByUserID", uint(1)).Return(orders, nil)

	result, err := svc.GetOrdersByUser(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}
