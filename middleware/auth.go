package middleware

import (
	"go-bookstore/dto"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &gin.H{"message": "Request must contain authorization token"})
			return
		}

		claims := dto.JWTClaim{}
		_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRECT")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &gin.H{"error": err.Error()})
			return
		}

		c.Set("user", &dto.AuthorizedUser{UserId: claims.UserId, Role: claims.Role})
		c.Next()
	}
}
