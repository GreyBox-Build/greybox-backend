package middlewares

import (
	"net/http"

	"slices"

	"github.com/gin-gonic/gin"
)

func AllowedHosts(allowedHosts []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestHost := c.Request.Host
		if slices.Contains(allowedHosts, requestHost) {
			return
		}

		// If the host is not allowed, return a 403 Forbidden response
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: oops!!! Host not allowed"})
		c.Abort()
	}
}
