package handlers

import (
	"errors"
	"go-bookstore/dto"
	"go-bookstore/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) AuthHandler {
	return AuthHandler{db: db}
}

func generateToken(userId uint, role models.UserRole) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRECT"))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["userId"] = userId
	claims["role"] = role

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthHandler) Login(c *gin.Context) {
	var body dto.LoginRequest
	var user models.User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := a.db.Where("email = ?", body.Email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, &gin.H{"message": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, &gin.H{"error": result.Error.Error()})
		return
	}

	if !checkPasswordHash(body.Password, user.Password) {
		c.JSON(http.StatusNotFound, &gin.H{"message": "Invalid email or password"})
		return
	}

	token, err := generateToken(user.ID, user.Role)

	if err != nil {
		c.JSON(http.StatusNotFound, &gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &dto.LoginResponse{Token: token})

}

func (a *AuthHandler) Register(c *gin.Context) {
	var body dto.RegisterRequest
	var user models.User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := a.db.Where("email = ?", body.Email).First(&user)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected >= 1 {
		c.JSON(http.StatusConflict, &gin.H{"message": "Email already exists"})
		return
	}

	hashedPass, err := HashPassword(body.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
		return
	}

	newUser := &models.User{Name: body.Name, Email: body.Email, Password: hashedPass}
	if createResult := a.db.Create(&newUser); createResult.Error != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"error": createResult.Error.Error()})
	}

	c.JSON(http.StatusCreated, &dto.RegisterResponse{Id: int(newUser.ID), Name: newUser.Name, Email: newUser.Email})
}
