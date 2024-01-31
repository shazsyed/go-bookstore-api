package handlers

import (
	"errors"
	"go-bookstore/dto"
	"go-bookstore/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookHandler struct {
	db *gorm.DB
}

func NewBookHandler(db *gorm.DB) BookHandler {
	return BookHandler{db: db}
}

func (b *BookHandler) GetAllBooks(c *gin.Context) {
	var books []dto.BookResponse
	b.db.Model(&models.Book{}).Find(&books)

	c.JSON(http.StatusOK, &gin.H{"data": books})
}

func (b *BookHandler) GetBookById(c *gin.Context) {
	var book []dto.BookResponse
	bookId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = b.db.Model(&models.Book{}).First(&book, bookId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, &gin.H{"message": "Not Found"})
		return
	}

	c.JSON(http.StatusOK, &gin.H{"data": book})
}

func (b *BookHandler) CreateBook(c *gin.Context) {
	var body dto.CreateBookRequest
	var response dto.BookResponse

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := &models.Book{Name: body.Name, Author: body.Author, Genre: body.Genre, ISBN: body.ISBN, Description: body.Description, PublishedYear: body.PublishedYear, Stock: body.Stock, Price: body.Price}
	result := b.db.Create(book)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": result.Error.Error()})
		return
	}

	response = dto.BookResponse{ID: book.ID, Name: book.Name, Author: book.Author, Genre: book.Genre, ISBN: book.ISBN, Description: book.Description, PublishedYear: book.PublishedYear, Stock: book.Stock, Price: book.Price}
	c.JSON(http.StatusCreated, &gin.H{"data": response})
}

func (b *BookHandler) UpdateBook(c *gin.Context) {
	var body dto.UpdateBookRequest
	var book models.Book

	bookId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := b.db.Model(&models.Book{}).First(&book, bookId); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, &gin.H{"message": "Not Found"})
			return
		}
	}

	if result := b.db.Model(&book).Updates(&models.Book{Name: body.Name, Author: body.Author, Description: body.Description, ISBN: body.ISBN, Genre: body.Genre, Stock: body.Stock, Price: body.Price}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": result.Error.Error()})
		return
	}

	response := dto.BookResponse{ID: book.ID, Name: book.Name, Author: book.Author, Genre: book.Genre, ISBN: book.ISBN, Description: book.Description, PublishedYear: book.PublishedYear, Stock: book.Stock, Price: book.Price}
	c.JSON(http.StatusOK, &response)

}

func (b *BookHandler) DeleteBook(c *gin.Context) {
	bookId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := b.db.Delete(&models.Book{}, bookId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, &gin.H{"message": "Not found"})
		return
	}

	c.JSON(http.StatusOK, &gin.H{"message": "success"})
}
