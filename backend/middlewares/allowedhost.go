package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AllowedHosts(allowedHosts []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestHost := c.Request.Host
		for _, allowedHost := range allowedHosts {
			if requestHost == allowedHost {
				return
			}
		}

		// If the host is not allowed, return a 403 Forbidden response
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: oops!!! Host not allowed"})
		c.Abort()
	}
}
