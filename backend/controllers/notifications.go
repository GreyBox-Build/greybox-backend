package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils"
	"fmt"
	"net/http"
	"time"

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

func TatumWebhook(c *gin.Context) {
	var user models.User
	var err error
	var masterWallet models.MasterWallet
	var currentBalance float32
	resultChan := make(chan string)
	errChan := make(chan error)

	authHeader := c.GetHeader("x-payload-hash")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var input serializers.Webhook
	if c.ShouldBindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	check := models.VerifyWebhookAuthenticity(authHeader, input)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	go apis.GetExchangeRate("USD", input.Chain, resultChan, errChan)
	if input.Type == "native" {
		c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
		return
	}
	user, err = models.FindUserByAddress(input.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	masterWallet, err = models.FetchMasterWallet(user.CryptoCurrency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	currentBalance, err = apis.FetchAccountBalanceCUSD(input.Address)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transType := ""
	select {
	case exchangeRate := <-resultChan:
		txId := input.TxID
		if currentBalance > user.PreviousBalance {
			user.PreviousBalance = currentBalance
			transType = "Deposit"
			amount := utils.PerformAmountCalculation(exchangeRate, input.Amount)

			_, code, err := apis.PerformTransactionCelo(amount, user.AccountAddress, masterWallet.PrivateKey, true)
			if err != nil {
				c.JSON(code, gin.H{"error": err.Error()})
				return
			}

		} else {
			transType = "Withdrawal"
		}
		var trans models.Transaction
		trans.UserID = user.ID
		trans.User = user
		trans.Hash = txId
		trans.TransactionId = txId
		trans.Amount = input.Amount
		trans.Description = fmt.Sprintf("%s of %s %s", transType, input.Amount, input.Asset)
		trans.Address = input.Address
		trans.BlockNumber = uint(input.BlockNumber)
		trans.CounterAddress = input.CounterAddress
		trans.Chain = input.Chain
		trans.TransactionSubType = transType
		trans.TransactionType = input.Type
		trans.Status = "Completed"
		if err := trans.SaveTransaction(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.UpdateUser()
		c.JSON(http.StatusOK, gin.H{"errors": false, "status": "webhook received successfully"})
		return
	case err := <-errChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	case <-time.After(10 * time.Second): // Set a timeout for the request
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
		return
	}

}
