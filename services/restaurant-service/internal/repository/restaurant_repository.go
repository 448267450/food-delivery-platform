package repository

import (
	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/model"
	"gorm.io/gorm"
)

// RestaurantRepository defines all database operations for restaurants
type RestaurantRepository interface {
	Create(restaurant *model.Restaurant) error
	FindAll() ([]model.Restaurant, error)
	FindByID(id uint) (*model.Restaurant, error)
	FindByOwnerID(ownerID uint) ([]model.Restaurant, error)
	Update(restaurant *model.Restaurant) error
	Delete(id uint) error
}

type restaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}

func (r *restaurantRepository) Create(restaurant *model.Restaurant) error {
	return r.db.Create(restaurant).Error
}

// FindAll returns all restaurants without menu items (for listing page)
func (r *restaurantRepository) FindAll() ([]model.Restaurant, error) {
	var restaurants []model.Restaurant
	err := r.db.Find(&restaurants).Error
	return restaurants, err
}

// FindByID returns a single restaurant WITH its menu items
func (r *restaurantRepository) FindByID(id uint) (*model.Restaurant, error) {
	var restaurant model.Restaurant
	err := r.db.Preload("MenuItems").First(&restaurant, id).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *restaurantRepository) FindByOwnerID(ownerID uint) ([]model.Restaurant, error) {
	var restaurants []model.Restaurant
	err := r.db.Where("owner_id = ?", ownerID).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) Update(restaurant *model.Restaurant) error {
	return r.db.Save(restaurant).Error
}

func (r *restaurantRepository) Delete(id uint) error {
	return r.db.Delete(&model.Restaurant{}, id).Error
}

// MenuItemRepository defines all database operations for menu items
type MenuItemRepository interface {
	Create(item *model.MenuItem) error
	FindByID(id uint) (*model.MenuItem, error)
	FindByRestaurantID(restaurantID uint) ([]model.MenuItem, error)
	Update(item *model.MenuItem) error
	Delete(id uint) error
}

type menuItemRepository struct {
	db *gorm.DB
}

func NewMenuItemRepository(db *gorm.DB) MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) Create(item *model.MenuItem) error {
	return r.db.Create(item).Error
}

func (r *menuItemRepository) FindByID(id uint) (*model.MenuItem, error) {
	var item model.MenuItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *menuItemRepository) FindByRestaurantID(restaurantID uint) ([]model.MenuItem, error) {
	var items []model.MenuItem
	err := r.db.Where("restaurant_id = ?", restaurantID).Find(&items).Error
	return items, err
}

func (r *menuItemRepository) Update(item *model.MenuItem) error {
	return r.db.Save(item).Error
}

func (r *menuItemRepository) Delete(id uint) error {
	return r.db.Delete(&model.MenuItem{}, id).Error
}
