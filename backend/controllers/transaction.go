package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils"
	"backend/utils/signing"
	"backend/utils/tokens"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	chain := c.Query("chain")
	transactions, err := models.GetTransactionsByUserID(userId, strings.ToUpper(chain))
	if err != nil {
		c.JSON(400, gin.H{
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
	transaction, err := models.GetTransactionByHash(hash, strings.ToUpper(chain))
	if err != nil {
		c.JSON(400, gin.H{
			"status": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"errors": false,
		"data":   transaction,
		"status": "retrieved transaction successfully",
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
	amount, accountAddress, Chain := input.Amount, input.AccountAddress, input.Chain
	switch strings.ToUpper(Chain) {
	case serializers.Chains.Celo:

		txHash, code, err := apis.PerformTransactionCelo(amount, accountAddress, user.PrivateKey, false)
		if err != nil {
			c.JSON(code, gin.H{"error": err.Error(), "message": "transaction failed"})
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
		hash := txData["txId"]

		data := map[string]interface{}{
			"transaction_hash": hash,
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

func AmountToReceive(c *gin.Context) {
	amount := c.Query("amount")
	currency := c.Query("currency")
	asset := c.Query("cryptoAsset")
	transType := c.Query("type")
	rateChan := make(chan string)
	errChan := make(chan error)
	go apis.GetExchangeRate(currency, strings.ToUpper(asset), rateChan, errChan)
	percent_reduction := utils.CalculateOnePercent(amount)
	amountConv, _ := strconv.ParseFloat(amount, 32)
	percentConv, _ := strconv.ParseFloat(percent_reduction, 32)
	newAmount := amountConv - percentConv
	select {
	case rate := <-rateChan:
		AssetAmount := ""
		data := map[string]string{
			"asset": strings.ToUpper(asset),
		}
		switch transType {
		case "on-ramp":
			AssetAmount = utils.ConvertTokenToNative(rate, strconv.FormatFloat(newAmount, 'f', 2, 64))
		case "off-ramp":
			AssetAmount = utils.ConvertAssetToFiat(rate, strconv.FormatFloat(newAmount, 'f', 2, 64))
			data["asset"] = strings.ToUpper(currency)
		}

		data["amount"] = AssetAmount

		c.JSON(200, gin.H{
			"errors": false,
			"status": "calculated amount to receive",
			"data":   data,
		})
		return

	case err := <-errChan:
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	case <-time.After(10 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
		return
	}

}
