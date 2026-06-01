package repository

import (
	"github.com/448267450/food-delivery-platform/services/order-service/internal/model"
	"gorm.io/gorm"
)

// OrderRepository defines all database operations for orders
type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id uint) (*model.Order, error)
	FindByUserID(userID uint) ([]model.Order, error)
	FindByRestaurantID(restaurantID uint) ([]model.Order, error)
	UpdateStatus(id uint, status model.OrderStatus) error
	Update(order *model.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create inserts a new order with its items in a single transaction
func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert the order first
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindByID returns an order WITH its items preloaded
func (r *orderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByUserID returns all orders for a given user, newest first
func (r *orderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// FindByRestaurantID returns all orders for a given restaurant, newest first
func (r *orderRepository) FindByRestaurantID(restaurantID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Items").
		Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// UpdateStatus only updates the status field, not the whole record
func (r *orderRepository) UpdateStatus(id uint, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Update saves the full order record
func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}
