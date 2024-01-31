package dto

import (
	"go-bookstore/models"

	"github.com/golang-jwt/jwt"
)

type LoginRequest struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type RegisterResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"  binding:"required,email"`
}

type AuthorizedUser struct {
	UserId uint `json:"userId"`
	Role   models.UserRole `json:"role"`
}

type JWTClaim struct {
	UserId uint            `json:"userId"`
	Role   models.UserRole `json:"role"`
	jwt.StandardClaims
}
