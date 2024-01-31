package models

import (
	"gorm.io/gorm"
)

type UserRole string

const (
	CUSTOMER_ROLE UserRole = "customer"
	ADMIN_ROLE    UserRole = "admin"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Password string `gorm:"not null"`
	Cart     Cart   `gorm:"foreignKey:UserID"`
	// Role     UserRole `gorm:"type:enum('admin','customer');default:'customer'"`
	Role UserRole `gorm:"default:'customer'"`
}
