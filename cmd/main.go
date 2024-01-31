package main

import (
	"go-bookstore/database"
	"go-bookstore/handlers"
	"go-bookstore/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := database.ConnectDB(false) //FLAG USED TO LOAD MOCK DATA IN DATABASE OR NOT FOR TESTING, CHANGE IT TO FALSE AFTER 1st RUN
	r := gin.Default()

	booksHandler := handlers.NewBookHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	cartHandler := handlers.NewCartHandler(db)

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/login", authHandler.Login)
		authRoute.POST("/register", authHandler.Register)
	}

	booksRoute := r.Group("/book")
	booksRoute.Use(middleware.AuthMiddleware())
	{
		booksRoute.GET("/all", booksHandler.GetAllBooks)
		booksRoute.GET("/:id", booksHandler.GetBookById)
		booksRoute.POST("/add", middleware.MustBeAdmin(), booksHandler.CreateBook)
		booksRoute.PUT("/:id", middleware.MustBeAdmin(), booksHandler.UpdateBook)
		booksRoute.DELETE("/:id", middleware.MustBeAdmin(), booksHandler.DeleteBook)
	}

	cartRoute := r.Group("/cart")
	cartRoute.Use(middleware.AuthMiddleware())
	{
		cartRoute.GET("/", cartHandler.GetCart)
		cartRoute.DELETE("/item/:id", cartHandler.DeleteBookInCart)
		cartRoute.POST("/add", cartHandler.AddBookInCart)
	}

	r.Run()
}
