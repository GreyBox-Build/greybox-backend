package middlewares

import (
	"backend/utils/tokens"
	"net/http"

	"backend/models"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := tokens.ExtractToken(c)

		flag := tokens.IsTokenValid(tokenString)

		if !flag {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := tokens.ExtractUserID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		user, err := models.GetUserByID(id)
		if err != nil {
			c.JSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}
		if user.Role != "Admin" {
			c.AbortWithStatusJSON(403, gin.H{"error": "You dont have permission to perform this action"})
			c.Abort()
			return
		}
		c.Next()
	}
}
