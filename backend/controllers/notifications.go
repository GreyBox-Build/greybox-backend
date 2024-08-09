package controllers

import (
	"backend/apis"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHmac(c *gin.Context) {
	err := apis.RegisterHmac()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"message": "HMAC registered successfully"})
	c.JSON(http.StatusOK, gin.H{"message": "HMAC registered successfully"})
}
