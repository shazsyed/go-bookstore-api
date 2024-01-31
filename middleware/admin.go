package middleware

import (
	"go-bookstore/dto"
	"go-bookstore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MustBeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")

		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &gin.H{"message": "Unauthorized"})
			return
		}
		
		typedUser := user.(*dto.AuthorizedUser)

		if typedUser.Role != models.ADMIN_ROLE {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &gin.H{"message": "Must be an admin to perform this action"})
			return
		}

		c.Next()

	}
}
