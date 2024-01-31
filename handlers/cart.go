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

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) CartHandler {
	return CartHandler{db: db}
}

func generateCartResponse(cartItems *[]models.CartItem) *dto.CartResponse {
	var response dto.CartResponse

	var total float64
	for _, item := range *cartItems {
		total += item.Book.Price * float64(item.Quantity)
		response.Items = append(response.Items, dto.CartItem{Book: dto.BookResponse{
			ID:            item.Book.ID,
			Name:          item.Book.Name,
			Author:        item.Book.Author,
			Genre:         item.Book.Genre,
			Description:   item.Book.Description,
			ISBN:          item.Book.ISBN,
			PublishedYear: item.Book.PublishedYear,
			Price:         item.Book.Price,
			Stock:         item.Book.Stock,
		},
			Quantity: item.Quantity,
			SubTotal: item.Book.Price * float64(item.Quantity),
		})
	}
	response.Total = total

	return &response
}

func (ch *CartHandler) AddBookInCart(c *gin.Context) {
	var user models.User
	var book models.Book
	var cart models.Cart
	var body dto.AddCartItemRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqUser, exists := c.Get("user")
	typedReqUser := reqUser.(*dto.AuthorizedUser)

	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{"message": "Unexpected error, user not found in the context"})
		return
	}

	ch.db.First(&user, typedReqUser.UserId)
	if result := ch.db.First(&book, body.BookId); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, &gin.H{"message": "No book found with such id"})
			return
		}
		c.JSON(http.StatusInternalServerError, &gin.H{"error": result.Error.Error()})
		return
	}

	if body.Quantity > book.Stock {
		c.JSON(http.StatusOK, &gin.H{"message": "Not enough stock available for this book"})
		return
	}

	if err := ch.db.Model(&user).Preload("CartItem").Association("Cart").Find(&cart); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	cart.UserID = user.ID
	var bookAlreadyInCart bool = false

	for i, v := range cart.CartItem {
		if v.BookID == body.BookId {
			bookAlreadyInCart = true
			cart.CartItem[i].Quantity = cart.CartItem[i].Quantity + body.Quantity
			break
		}
	}

	if !bookAlreadyInCart {
		cartItem := models.CartItem{Book: book, Quantity: body.Quantity}
		cart.CartItem = append(cart.CartItem, cartItem)
	}

	tx := ch.db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if err := tx.Transaction(func(tx *gorm.DB) error {

		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&cart).Error; err != nil {
			return err
		}

		if err := tx.Model(&book).Update("stock", book.Stock-body.Quantity).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	ch.db.Preload("Cart.CartItem.Book").First(&user, typedReqUser.UserId)

	response := generateCartResponse(&user.Cart.CartItem)
	c.JSON(http.StatusCreated, &gin.H{"data": response})
}

func (ch *CartHandler) GetCart(c *gin.Context) {
	var user models.User
	reqUser, exists := c.Get("user")
	typedReqUser := reqUser.(*dto.AuthorizedUser)

	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{"message": "Unexpected error, user not found in the context"})
		return
	}

	ch.db.Preload("Cart.CartItem.Book").First(&user, typedReqUser.UserId)

	response := generateCartResponse(&user.Cart.CartItem)
	c.JSON(http.StatusOK, &gin.H{"data": response})
}

func (ch *CartHandler) DeleteBookInCart(c *gin.Context) {
	var user models.User
	var item models.CartItem
	var book models.Book

	bookId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqUser, exists := c.Get("user")
	typedReqUser := reqUser.(*dto.AuthorizedUser)

	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{"message": "Unexpected error, user not found in the context"})
		return
	}

	ch.db.Preload("Cart.CartItem.Book").First(&user, typedReqUser.UserId)

	if err := ch.db.Model(&user.Cart.CartItem).Where("book_id = ?", bookId).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, &gin.H{"message": "No book found in the cart with such id"})
			return
		}
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	tx := ch.db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if err := tx.Transaction(func(tx *gorm.DB) error {
		if item.Quantity <= 1 {
			if err := ch.db.Delete(&item).Error; err != nil {
				return err
			}
		} else if err := ch.db.Model(&item).Update("quantity", item.Quantity-1).Error; err != nil {
			return err
		}

		if err := ch.db.First(&book, bookId).Error; err != nil {
			return err
		}

		if err := ch.db.Model(&book).Update("stock", book.Stock+1).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	ch.db.Preload("Cart.CartItem.Book").First(&user, typedReqUser.UserId)

	response := generateCartResponse(&user.Cart.CartItem)
	c.JSON(http.StatusOK, &gin.H{"data": response})

}
