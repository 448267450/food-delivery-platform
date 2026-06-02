package service_test

import (
	"errors"
	"testing"

	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/model"
	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==================== Mock Repositories ====================

type MockRestaurantRepository struct {
	mock.Mock
}

func (m *MockRestaurantRepository) Create(r *model.Restaurant) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockRestaurantRepository) FindAll() ([]model.Restaurant, error) {
	args := m.Called()
	return args.Get(0).([]model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) FindByID(id uint) (*model.Restaurant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) FindByOwnerID(ownerID uint) ([]model.Restaurant, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]model.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) Update(r *model.Restaurant) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockRestaurantRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockMenuItemRepository struct {
	mock.Mock
}

func (m *MockMenuItemRepository) Create(item *model.MenuItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockMenuItemRepository) FindByID(id uint) (*model.MenuItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MenuItem), args.Error(1)
}

func (m *MockMenuItemRepository) FindByRestaurantID(restaurantID uint) ([]model.MenuItem, error) {
	args := m.Called(restaurantID)
	return args.Get(0).([]model.MenuItem), args.Error(1)
}

func (m *MockMenuItemRepository) Update(item *model.MenuItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockMenuItemRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// ==================== Restaurant Tests ====================

func TestCreateRestaurant_Success(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	req := &service.CreateRestaurantRequest{
		OwnerID:     1,
		Name:        "Ryan's Burger",
		Description: "Best burgers in Austin",
		Address:     "123 Main St",
		Phone:       "512-000-1234",
	}

	mockRestRepo.On("Create", mock.AnythingOfType("*model.Restaurant")).Return(nil)

	restaurant, err := svc.CreateRestaurant(req)

	assert.NoError(t, err)
	assert.NotNil(t, restaurant)
	assert.Equal(t, req.Name, restaurant.Name)
	assert.Equal(t, req.OwnerID, restaurant.OwnerID)
	assert.True(t, restaurant.IsOpen) // default should be open
	mockRestRepo.AssertExpectations(t)
}

func TestGetRestaurantByID_Success(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	restaurant := &model.Restaurant{
		ID:      1,
		OwnerID: 1,
		Name:    "Ryan's Burger",
		IsOpen:  true,
	}

	mockRestRepo.On("FindByID", uint(1)).Return(restaurant, nil)

	result, err := svc.GetRestaurantByID(1)

	assert.NoError(t, err)
	assert.Equal(t, restaurant.Name, result.Name)
}

func TestGetRestaurantByID_NotFound(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	mockRestRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetRestaurantByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "restaurant not found", err.Error())
}

func TestUpdateRestaurant_Success(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	existing := &model.Restaurant{
		ID:     1,
		Name:   "Old Name",
		IsOpen: true,
	}

	isOpen := false
	req := &service.UpdateRestaurantRequest{
		Name:   "New Name",
		IsOpen: &isOpen,
	}

	mockRestRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRestRepo.On("Update", mock.AnythingOfType("*model.Restaurant")).Return(nil)

	result, err := svc.UpdateRestaurant(1, req)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
	assert.False(t, result.IsOpen)
}

func TestDeleteRestaurant_NotFound(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	mockRestRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.DeleteRestaurant(999)

	assert.Error(t, err)
	assert.Equal(t, "restaurant not found", err.Error())
	mockRestRepo.AssertNotCalled(t, "Delete")
}

// ==================== Menu Item Tests ====================

func TestAddMenuItem_Success(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	restaurant := &model.Restaurant{ID: 1, Name: "Ryan's Burger"}
	req := &service.AddMenuItemRequest{
		Name:     "Classic Cheeseburger",
		Price:    12.99,
		Category: "burger",
	}

	mockRestRepo.On("FindByID", uint(1)).Return(restaurant, nil)
	mockMenuRepo.On("Create", mock.AnythingOfType("*model.MenuItem")).Return(nil)

	item, err := svc.AddMenuItem(1, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Name, item.Name)
	assert.Equal(t, req.Price, item.Price)
	assert.Equal(t, uint(1), item.RestaurantID)
	assert.True(t, item.IsAvailable)
}

func TestAddMenuItem_RestaurantNotFound(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	mockRestRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	item, err := svc.AddMenuItem(999, &service.AddMenuItemRequest{
		Name:  "Burger",
		Price: 10.00,
	})

	assert.Error(t, err)
	assert.Nil(t, item)
	assert.Equal(t, "restaurant not found", err.Error())
	mockMenuRepo.AssertNotCalled(t, "Create")
}

func TestUpdateMenuItem_WrongRestaurant(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	// Item belongs to restaurant 2, but we're trying to update via restaurant 1
	item := &model.MenuItem{
		ID:           1,
		RestaurantID: 2,
		Name:         "Burger",
		Price:        10.00,
	}

	mockMenuRepo.On("FindByID", uint(1)).Return(item, nil)

	price := 15.00
	result, err := svc.UpdateMenuItem(1, 1, &service.UpdateMenuItemRequest{
		Price: &price,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "menu item does not belong to this restaurant", err.Error())
	mockMenuRepo.AssertNotCalled(t, "Update")
}

func TestDeleteMenuItem_WrongRestaurant(t *testing.T) {
	mockRestRepo := new(MockRestaurantRepository)
	mockMenuRepo := new(MockMenuItemRepository)
	svc := service.NewRestaurantService(mockRestRepo, mockMenuRepo)

	item := &model.MenuItem{
		ID:           1,
		RestaurantID: 2, // belongs to restaurant 2
	}

	mockMenuRepo.On("FindByID", uint(1)).Return(item, nil)

	err := svc.DeleteMenuItem(1, 1) // trying via restaurant 1

	assert.Error(t, err)
	assert.Equal(t, "menu item does not belong to this restaurant", err.Error())
	mockMenuRepo.AssertNotCalled(t, "Delete")
}
