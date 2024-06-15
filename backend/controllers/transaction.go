package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/tokens"
	"math/big"
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
	transactions, err := apis.GetUserTransactions(strings.ToLower(user.CryptoCurrency), user.AccountAddress, "zero-transfer", 50)
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

	switch strings.ToLower(Chain) {
	case "celo":
		resultChan := make(chan map[string]interface{})
		errChan := make(chan error)

		go func() {
			result, err := apis.CalculateEstimatedFeeCelo(amount, accountAddress, user.AccountAddress)
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			// Handle the result from CalculateEstimatedFeeCelo
			gprice, _ := result["gasPrice"].(string)
			price, err := strconv.ParseInt(gprice, 10, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "conversion failed for conversion of gas price"})
				return
			}
			gasPrice := models.WeiToGwei(big.NewInt(price))
			trans.TransFee, _ = result["gasLimit"].(float64)
			txHash, err := apis.PerformTransactionCelo(amount, accountAddress, user.PrivateKey, gasPrice.String(), result["gasLimit"].(float64))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "transaction failed"})
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

		case err := <-errChan:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "gas price could not be calculated"})
			return
		}

	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid chain"})

}
