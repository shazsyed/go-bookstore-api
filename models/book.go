package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Author        string `gorm:"not null"`
	Genre         string
	PublishedYear uint   `gorm:"not null;default:0" json:"published_year"`
	ISBN          string `gorm:"not null"`
	Description   string
	Stock         uint    `gorm:"not null;default:1"`
	Price         float64 `gorm:"not null;default:1"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	if b.PublishedYear == 0 {
		b.PublishedYear = uint(time.Now().Year())
	}
	return
}
