package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/tokens"
	"net/http"
	"os"
	"strconv"
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
	category := c.Query("category")
	pageSize := c.Query("pageSize")
	page, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"errors": true,
		})
		return
	}
	if len(category) != 0 || page != 0 {
		transactions, err := apis.GetUserTransactions(strings.ToLower(user.CryptoCurrency), user.AccountAddress, category, uint64(page))
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"errors": false,
			"data":   transactions,
			"status": "retrieved transactions successfully",
		})
		return

	}
	transactions, err := apis.GetUserTransactions(strings.ToLower(user.CryptoCurrency), user.AccountAddress, "", 50)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"errors": false,
		"data":   transactions,
		"status": "retrieved transactions successfully",
	})

}

func GetTransactionsByHash(c *gin.Context) {
	hash := c.Query("hash")
	chain := c.Query("chain")
	transactions, err := apis.GetTransactionByHash(chain, hash)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"errors": false,
		"data":   transactions,
		"status": "retrieved transactions successfully",
	})
}

func OffRampTransaction(c *gin.Context) {
	var input serializers.OffRampForm
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

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
	var trans models.Transaction
	trans.UserID = userId
	trans.Amount = input.Amount
	trans.User = user
	trans.Chain = input.Chain
	trans.Address = input.AccountAddress
	trans.Status = "pending"
	trans.TransactionSubType = "outgoing"

	amount, accountAddress, Chain := input.Amount, input.AccountAddress, input.Chain

	switch strings.ToUpper(Chain) {
	case serializers.Chains.Celo:

		txHash, code, err := apis.PerformTransactionCelo(amount, accountAddress, user.PrivateKey)
		if err != nil {
			c.JSON(code, gin.H{"error": err.Error(), "message": "transaction failed"})
			return
		}
		trans.Hash = txHash
		trans.Status = "completed"
		if err := trans.SaveTransaction(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "saving transaction info failed"})
			return
		}
		data := map[string]interface{}{
			"transaction_hash": txHash,
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"errors": false,
				"data":   data,
				"status": "transaction perform successfully",
			},
		)
		return
	case serializers.Chains.Stellar:
		transferData := serializers.TransferXLM{
			Amount:        amount,
			To:            accountAddress,
			FromSecret:    user.PrivateKey,
			Initialize:    false,
			Token:         "USDC",
			IssuerAccount: user.AccountAddress,
		}
		txData, code, err := apis.PerformTransactionXLM(transferData)
		if err != nil {
			c.JSON(code, gin.H{"error": err.Error(), "message": "transaction failed"})
			return
		}
		trans.Hash = txData["txId"]
		trans.Status = "completed"
		if err := trans.SaveTransaction(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "saving transaction into db failed"})
			return
		}
		data := map[string]interface{}{
			"transaction_hash": trans.Hash,
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"errors": false,
				"data":   data,
				"status": "transaction perform successfully",
			},
		)
		return

	}

}
