package model

import "time"

// OrderStatus represents the state machine for an order
// Flow: PENDING -> PAID -> PREPARING -> OUT_FOR_DELIVERY -> DELIVERED
//                       -> CANCELLED (from PENDING or PAID only)
type OrderStatus string

const (
	StatusPending         OrderStatus = "PENDING"          // order created, awaiting payment
	StatusPaid            OrderStatus = "PAID"             // payment confirmed
	StatusPreparing       OrderStatus = "PREPARING"        // restaurant is preparing
	StatusOutForDelivery  OrderStatus = "OUT_FOR_DELIVERY" // driver picked up
	StatusDelivered       OrderStatus = "DELIVERED"        // successfully delivered
	StatusCancelled       OrderStatus = "CANCELLED"        // cancelled
)

// ValidTransitions defines allowed state transitions
var ValidTransitions = map[OrderStatus][]OrderStatus{
	StatusPending:        {StatusPaid, StatusCancelled},
	StatusPaid:           {StatusPreparing, StatusCancelled},
	StatusPreparing:      {StatusOutForDelivery},
	StatusOutForDelivery: {StatusDelivered},
	StatusDelivered:      {},
	StatusCancelled:      {},
}

func (o *Order) CanTransitionTo(next OrderStatus) bool {
	allowed := ValidTransitions[o.Status]
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

type Order struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint        `gorm:"not null" json:"user_id"`
	RestaurantID uint        `gorm:"not null" json:"restaurant_id"`
	DriverID     *uint       `json:"driver_id"`
	Status       OrderStatus `gorm:"default:PENDING" json:"status"`
	TotalPrice   float64     `gorm:"not null" json:"total_price"`
	Address      string      `gorm:"not null" json:"delivery_address"`
	Note         string      `json:"note"`
	Items        []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID         uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID    uint    `gorm:"not null" json:"order_id"`
	MenuItemID uint    `gorm:"not null" json:"menu_item_id"`
	Name       string  `gorm:"not null" json:"name"`  // snapshot at order time
	Price      float64 `gorm:"not null" json:"price"` // snapshot at order time
	Quantity   int     `gorm:"not null" json:"quantity"`
}
