package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Role      string `gorm:"type:varchar(20);default:'user'" json:"role"` // 'admin' or 'user'
	Orders    []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}
