package controllers

import (
	"backend/models"
	"backend/utils/tokens"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RetrieveOnRampParamsV1(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"errors": true,
		})
		return
	}
	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	data := map[string]interface{}{
		"wallet_address":       user.AccountAddress,
		"asset":                "CUSD",
		"email":                user.Email,
		"external_customer_id": user.CustomerId,
		"x-client-id":          os.Getenv("X_CLIENT_ID"),
		"x-client-secret":      os.Getenv("X_CLIENT_SECRET"),
		"network":              "CELO",
		"country":              user.Country,
		"source_param":         os.Getenv("SOURCE_PARAM"),
	}

	c.JSON(200, gin.H{
		"errors": false,
		"data":   data,
		"status": "retrieved on ramp data successfully",
	})

}

func GetUserTransactions(c *gin.Context) {

}
