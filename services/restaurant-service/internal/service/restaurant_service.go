package service

import (
	"errors"

	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/model"
	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/repository"
)

// ==================== Request / Response structs ====================

type CreateRestaurantRequest struct {
	OwnerID     uint   `json:"owner_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address" binding:"required"`
	Phone       string `json:"phone"`
	ImageURL    string `json:"image_url"`
}

type UpdateRestaurantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	ImageURL    string `json:"image_url"`
	IsOpen      *bool  `json:"is_open"` // pointer so we can distinguish false vs not provided
}

type AddMenuItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

type UpdateMenuItemRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"` // pointer so we can distinguish 0 vs not provided
	Category    string   `json:"category"`
	ImageURL    string   `json:"image_url"`
	IsAvailable *bool    `json:"is_available"`
}

// ==================== Interface ====================

type RestaurantService interface {
	CreateRestaurant(req *CreateRestaurantRequest) (*model.Restaurant, error)
	GetAllRestaurants() ([]model.Restaurant, error)
	GetRestaurantByID(id uint) (*model.Restaurant, error)
	GetRestaurantsByOwner(ownerID uint) ([]model.Restaurant, error)
	UpdateRestaurant(id uint, req *UpdateRestaurantRequest) (*model.Restaurant, error)
	DeleteRestaurant(id uint) error

	AddMenuItem(restaurantID uint, req *AddMenuItemRequest) (*model.MenuItem, error)
	UpdateMenuItem(restaurantID uint, itemID uint, req *UpdateMenuItemRequest) (*model.MenuItem, error)
	DeleteMenuItem(restaurantID uint, itemID uint) error
}

// ==================== Implementation ====================

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
	menuItemRepo   repository.MenuItemRepository
}

func NewRestaurantService(
	restaurantRepo repository.RestaurantRepository,
	menuItemRepo repository.MenuItemRepository,
) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
		menuItemRepo:   menuItemRepo,
	}
}

// -------------------- Restaurant --------------------

func (s *restaurantService) CreateRestaurant(req *CreateRestaurantRequest) (*model.Restaurant, error) {
	restaurant := &model.Restaurant{
		OwnerID:     req.OwnerID,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Phone:       req.Phone,
		ImageURL:    req.ImageURL,
		IsOpen:      true,
	}

	if err := s.restaurantRepo.Create(restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (s *restaurantService) GetAllRestaurants() ([]model.Restaurant, error) {
	return s.restaurantRepo.FindAll()
}

func (s *restaurantService) GetRestaurantByID(id uint) (*model.Restaurant, error) {
	restaurant, err := s.restaurantRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}
	return restaurant, nil
}

func (s *restaurantService) GetRestaurantsByOwner(ownerID uint) ([]model.Restaurant, error) {
	return s.restaurantRepo.FindByOwnerID(ownerID)
}

func (s *restaurantService) UpdateRestaurant(id uint, req *UpdateRestaurantRequest) (*model.Restaurant, error) {
	// First fetch the existing record
	restaurant, err := s.restaurantRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	// Only update fields that were provided
	if req.Name != "" {
		restaurant.Name = req.Name
	}
	if req.Description != "" {
		restaurant.Description = req.Description
	}
	if req.Address != "" {
		restaurant.Address = req.Address
	}
	if req.Phone != "" {
		restaurant.Phone = req.Phone
	}
	if req.ImageURL != "" {
		restaurant.ImageURL = req.ImageURL
	}
	if req.IsOpen != nil {
		restaurant.IsOpen = *req.IsOpen
	}

	if err := s.restaurantRepo.Update(restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (s *restaurantService) DeleteRestaurant(id uint) error {
	_, err := s.restaurantRepo.FindByID(id)
	if err != nil {
		return errors.New("restaurant not found")
	}
	return s.restaurantRepo.Delete(id)
}

// -------------------- Menu Items --------------------

func (s *restaurantService) AddMenuItem(restaurantID uint, req *AddMenuItemRequest) (*model.MenuItem, error) {
	// Verify restaurant exists before adding menu item
	_, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, errors.New("restaurant not found")
	}

	item := &model.MenuItem{
		RestaurantID: restaurantID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Category:     req.Category,
		ImageURL:     req.ImageURL,
		IsAvailable:  true,
	}

	if err := s.menuItemRepo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *restaurantService) UpdateMenuItem(restaurantID uint, itemID uint, req *UpdateMenuItemRequest) (*model.MenuItem, error) {
	item, err := s.menuItemRepo.FindByID(itemID)
	if err != nil {
		return nil, errors.New("menu item not found")
	}

	// Security check: make sure this item actually belongs to this restaurant
	if item.RestaurantID != restaurantID {
		return nil, errors.New("menu item does not belong to this restaurant")
	}

	// Only update fields that were provided
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.ImageURL != "" {
		item.ImageURL = req.ImageURL
	}
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}

	if err := s.menuItemRepo.Update(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *restaurantService) DeleteMenuItem(restaurantID uint, itemID uint) error {
	item, err := s.menuItemRepo.FindByID(itemID)
	if err != nil {
		return errors.New("menu item not found")
	}

	// Security check: make sure this item actually belongs to this restaurant
	if item.RestaurantID != restaurantID {
		return errors.New("menu item does not belong to this restaurant")
	}

	return s.menuItemRepo.Delete(itemID)
}
