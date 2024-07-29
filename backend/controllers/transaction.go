package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/signing"
	"backend/utils/tokens"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
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
		"external_customer_id": user.ID,
		"network":              "CELO",
		"country":              user.Country,
		"source_param":         os.Getenv("SOURCE_PARAM"),
	}
	switch user.CryptoCurrency {
	case serializers.Chains.Celo:
		data["asset"] = "CUSD"
		data["network"] = "CELO"
		data["x-client-id"] = os.Getenv("X_CLIENT_ID")
		data["x-client-secret"] = os.Getenv("X_CLIENT_SECRET")
	case serializers.Chains.Stellar:
		data["asset"] = "USDC"
		data["network"] = "XLM"
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
	if user.CryptoCurrency == "XLM" {
		if len(pageSize) == 0 {
			pageSize = "10"
		}
		trans, err := apis.GetUserTransactionXLM(user.AccountAddress, pageSize)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		decodTrans, err := apis.DecodeTransactionDataXLM(trans)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"errors": false,
			"data":   decodTrans,
			"status": "retrieved transactions successfully",
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
	switch strings.ToUpper(chain) {
	case serializers.Chains.Celo:
		transactions, err := apis.GetTransactionByHash(strings.ToLower(chain), hash)
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
	case serializers.Chains.Stellar:
		trans := []serializers.TransactionXLM{}
		transaction, err := apis.GetTransactionByHashXLM(hash)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		trans = append(trans, transaction)
		decodTrans, err := apis.DecodeTransactionDataXLM(trans)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"errors": false,
			"data":   decodTrans,
			"status": "retrieved transactions successfully",
		})
		return

	}

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
			FromAccount:   user.AccountAddress,
		}
		txData, code, err := apis.PerformTransactionXLM(transferData)
		if err != nil {
			c.JSON(code, gin.H{"error": err.Error(), "message": "transaction failed"})
			return
		}
		trans.Hash = txData["txId"]
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

func SignUrl(c *gin.Context) {
	var input serializers.SignUrl
	var err error
	var signedUrl string
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	signedUrl, err = signing.GenerateSignedURL(input.Url, os.Getenv("MOONPAY_API_TEST_KEY"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(signedUrl)
	c.JSON(200, gin.H{
		"signedUrl": signedUrl,
	})

}

func KMStransactionVerification(c *gin.Context) {
	transId := c.Param("transaction_id")
	trans, err := models.GetTransactionByHash(transId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	trans.Status = "completed"
	data := map[string]string{
		"transaction_hash": trans.Hash,
	}
	if err := trans.SaveTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"errors": false,
		"data":   data,
		"status": "retrieved transactions successfully",
	})

}

func FetchChain(c *gin.Context) {

	root, _ := os.Getwd()

	jsonFilePath := filepath.Join(root, "/templates", "/chains.json")

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []serializers.Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "status": "fetched accepted chains", "errors": false})
}
