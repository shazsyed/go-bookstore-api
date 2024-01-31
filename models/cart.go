package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	UserID   uint
	CartItem []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model

	CartID   uint
	BookID   uint
	Book     Book `gorm:"foreignKey:BookID"`
	Quantity uint `gorm:"default:1"`
}
