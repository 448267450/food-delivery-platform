package service

import (
	"errors"
	"fmt"

	"github.com/448267450/food-delivery-platform/services/order-service/internal/model"
	"github.com/448267450/food-delivery-platform/services/order-service/internal/repository"
)

// ==================== Request / Response structs ====================

type CreateOrderRequest struct {
	UserID       uint               `json:"user_id" binding:"required"`
	RestaurantID uint               `json:"restaurant_id" binding:"required"`
	Address      string             `json:"delivery_address" binding:"required"`
	Note         string             `json:"note"`
	Items        []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
	MenuItemID uint    `json:"menu_item_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Quantity   int     `json:"quantity" binding:"required,min=1"`
}

type UpdateStatusRequest struct {
	Status model.OrderStatus `json:"status" binding:"required"`
}

// ==================== Interface ====================

type OrderService interface {
	CreateOrder(req *CreateOrderRequest) (*model.Order, error)
	GetOrderByID(id uint) (*model.Order, error)
	GetOrdersByUser(userID uint) ([]model.Order, error)
	GetOrdersByRestaurant(restaurantID uint) ([]model.Order, error)
	UpdateOrderStatus(id uint, req *UpdateStatusRequest) (*model.Order, error)
	CancelOrder(id uint) error
}

// ==================== Implementation ====================

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) CreateOrder(req *CreateOrderRequest) (*model.Order, error) {
	// Build order items and calculate total price
	var items []model.OrderItem
	var totalPrice float64

	for _, i := range req.Items {
		item := model.OrderItem{
			MenuItemID: i.MenuItemID,
			Name:       i.Name,
			Price:      i.Price,
			Quantity:   i.Quantity,
		}
		items = append(items, item)
		totalPrice += i.Price * float64(i.Quantity)
	}

	order := &model.Order{
		UserID:       req.UserID,
		RestaurantID: req.RestaurantID,
		Address:      req.Address,
		Note:         req.Note,
		Status:       model.StatusPending,
		TotalPrice:   totalPrice,
		Items:        items,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *orderService) GetOrderByID(id uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *orderService) GetOrdersByUser(userID uint) ([]model.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *orderService) GetOrdersByRestaurant(restaurantID uint) ([]model.Order, error) {
	return s.orderRepo.FindByRestaurantID(restaurantID)
}

func (s *orderService) UpdateOrderStatus(id uint, req *UpdateStatusRequest) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// *** State machine validation ***
	// Check if the transition is allowed before writing to DB
	if !order.CanTransitionTo(req.Status) {
		return nil, fmt.Errorf(
			"invalid status transition: %s → %s",
			order.Status,
			req.Status,
		)
	}

	if err := s.orderRepo.UpdateStatus(id, req.Status); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Return the updated order
	order.Status = req.Status
	return order, nil
}

func (s *orderService) CancelOrder(id uint) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	// Can only cancel from PENDING or PAID
	if !order.CanTransitionTo(model.StatusCancelled) {
		return fmt.Errorf(
			"cannot cancel order with status: %s (only PENDING or PAID orders can be cancelled)",
			order.Status,
		)
	}

	return s.orderRepo.UpdateStatus(id, model.StatusCancelled)
}
