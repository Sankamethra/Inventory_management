package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint        `gorm:"not null" json:"user_id"`
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status     string      `gorm:"type:varchar(20);default:'pending'" json:"status"`
	TotalPrice float64     `gorm:"not null" json:"total_price"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"` // Price at the time of order
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
} 