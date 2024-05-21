package controllers

import (
	"backend/models"
	"backend/utils/tokens"
	"net/http"
	"os"
	"strings"

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
	trans, err := models.CreateTransaction(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	data := map[string]interface{}{
		"transaction_id":       trans.TransactionId,
		"wallet_address":       user.AccountAddress,
		"currency_code":        "cusd",
		"fiat_currency_code":   strings.ToLower(user.Currency),
		"email":                user.Email,
		"external_customer_id": user.CustomerId,
		"api_key":              os.Getenv("MOONPAY_API_TEST_KEY"),
	}

	c.JSON(200, gin.H{
		"errors": false,
		"data":   data,
		"status": "retrieved on ramp data successfully",
	})

}
