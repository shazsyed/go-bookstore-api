package database

import (
	"encoding/json"
	"errors"
	"go-bookstore/handlers"
	"go-bookstore/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func loadBookMockData(path string, db *gorm.DB) {

	type MockDataBooks struct {
		Books []models.Book `json:"books"`
	}

	var books MockDataBooks
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read books mock data, skipping inserting fake books in database, err: %v", err.Error())
		return
	}

	if err := json.Unmarshal(data, &books); err != nil {
		log.Printf("Failed to unmarshall mock books data, err: %v", err.Error())
		return
	}

	if err := db.CreateInBatches(books.Books, 100).Error; err != nil {
		log.Printf("Failed to insert book mock data in the database, err: %v", err.Error())
		return
	}

	log.Printf("Books mock data has been inserted")
}

func ConnectDB(loadFakeBooks bool) *gorm.DB {

	db, err := gorm.Open(sqlite.Open("./database/test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Database connection failed")
	}

	if err := db.AutoMigrate(&models.Book{}, &models.CartItem{}, &models.Cart{}, &models.User{}); err != nil {
		log.Fatalf("Failed to migrate, err %v", err.Error())
	}

	var adminUser models.User
	if result := db.First(&adminUser, "role = ?", models.ADMIN_ROLE); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			hash, err := handlers.HashPassword("admin")

			if err != nil {
				log.Println("Failed to create an admin account in the database")
				return db
			}

			db.Create(&models.User{
				Name:     "Admin",
				Email:    "admin@admin.com",
				Password: hash,
				Role:     models.ADMIN_ROLE,
				Cart:     models.Cart{},
			})
		}
	}

	if loadFakeBooks {
		loadBookMockData("./database/mock_data.json", db)
	}

	return db
}
