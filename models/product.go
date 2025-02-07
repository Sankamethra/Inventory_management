package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string        `gorm:"not null" json:"name"`
	Description   string        `json:"description"`
	Price         float64       `gorm:"not null" json:"price"`
	BasePrice     float64       `gorm:"not null" json:"base_price"`
	Stock         int           `gorm:"not null" json:"stock"`
	OrderItems    []OrderItem   `gorm:"foreignKey:ProductID" json:"order_items,omitempty"`
	PriceHistory []PriceHistory `gorm:"foreignKey:ProductID" json:"price_history,omitempty"`
}

type PriceHistory struct {
	gorm.Model
	ProductID uint    `gorm:"not null" json:"product_id"`
	OldPrice  float64 `gorm:"not null" json:"old_price"`
	NewPrice  float64 `gorm:"not null" json:"new_price"`
	Reason    string  `gorm:"type:varchar(100)" json:"reason"`
}