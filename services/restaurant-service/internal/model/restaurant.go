package model

import "time"

type Restaurant struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID     uint       `gorm:"not null" json:"owner_id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `json:"description"`
	Address     string     `gorm:"not null" json:"address"`
	Phone       string     `json:"phone"`
	ImageURL    string     `json:"image_url"`
	IsOpen      bool       `gorm:"default:true" json:"is_open"`
	MenuItems   []MenuItem `gorm:"foreignKey:RestaurantID" json:"menu_items,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type MenuItem struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RestaurantID uint      `gorm:"not null" json:"restaurant_id"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `json:"description"`
	Price        float64   `gorm:"not null" json:"price"`
	Category     string    `json:"category"` // e.g. "burger", "drink", "dessert"
	ImageURL     string    `json:"image_url"`
	IsAvailable  bool      `gorm:"default:true" json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
