package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"not null"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Status     string      `json:"status" gorm:"type:varchar(20);default:'pending'"`
	TotalPrice float64     `json:"total_price" gorm:"not null"`
	User       User        `json:"user" gorm:"foreignKey:UserID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}
