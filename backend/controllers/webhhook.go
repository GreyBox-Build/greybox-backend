package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func OnRampNotification(c *gin.Context) {
	var input serializers.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("error", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("onramp request body: ", input)

	// Fetch the request based on the EventObject ID
	request, err := models.GetHurupayRequestRequestId(input.EventObject.ID)
	if err != nil {
		fmt.Println("error", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	trans, err := createTransaction(request, input)
	if err != nil {
		fmt.Println("error", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Fetch the master wallet for the user's cryptocurrency
	masterWallet, err := models.FetchMasterWallet(request.User.CryptoCurrency)
	if err != nil {
		fmt.Println("error", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Process the event based on the event type
	if err = processEvent(input.EventType, *request, *trans, masterWallet); err != nil {
		fmt.Println("error", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"errors": false, "status": "processed notification successfully"})
}

// Create a transaction from the request and event input
func createTransaction(request *models.HurupayRequest, input serializers.Event) (*models.Transaction, error) {
	var trans models.Transaction
	previousTrans, err := models.GetTransactionByRequestId(request.RequestId)
	if err == nil && previousTrans != nil {
		hash := utils.LastPart(input.EventObject.BlockchainProof, "/")
		previousTrans.Hash = hash
		previousTrans.TransactionId = hash
		previousTrans.Address = request.User.AccountAddress
		previousTrans.TransactionSubType = "Deposit"
		if request.RequestType == models.OffRamp {
			previousTrans.TransactionSubType = "Withdrawal"
		}
		previousTrans.Asset = strings.ToUpper(input.EventObject.BlockchainToken)
		previousTrans.Chain = input.EventObject.BlockchainNetwork
		previousTrans.TransactionType = "Fungible Token"
		previousTrans.Amount = input.EventObject.TokenAmount.String()
		if request.RequestType == models.OffRamp {
			previousTrans.Amount = input.EventObject.CollectionAmount.String()
		}
		previousTrans.Status = "Pending"
		return previousTrans, nil
	}

	trans.UserID = request.User.ID
	trans.User = request.User
	trans.Description = input.EventObject.Description
	trans.RequestId = request.RequestId
	trans.TransactionSubType = "Deposit"
	if request.RequestType == models.OffRamp {
		trans.TransactionSubType = "Withdrawal"
	}

	if err := trans.SaveTransaction(); err != nil {
		return nil, err
	}
	return &trans, nil
}

// Process event based on event type
func processEvent(eventType string, request models.HurupayRequest, trans models.Transaction, masterWallet models.MasterWallet) error {
	switch eventType {
	case "collections.successful":
		request.Status = "Completed"
		trans.Status = "Completed"
		_ = trans.UpdateTransaction()
		_ = request.UpdateHurupayRequest()

		nativeAmount, err := utils.PerformDepositofNativeCalculation(trans.Amount, "USD", request.User.CryptoCurrency)
		if err != nil {
			return err
		}

		// Process native transaction based on the cryptocurrency
		return processNativeTransaction(request.User.CryptoCurrency, nativeAmount, request, masterWallet)

	default:
		request.Status = strings.ToUpper(utils.LastPart(eventType, "."))
		trans.Status = strings.ToUpper(utils.LastPart(eventType, "."))
		// Save the updated request and transaction
		_ = trans.UpdateTransaction()
		return request.UpdateHurupayRequest()
	}
}

// Handle native transaction processing for CELO and XLM
func processNativeTransaction(crypto string, nativeAmount string, request models.HurupayRequest, masterWallet models.MasterWallet) error {
	switch crypto {
	case "CELO":
		go processCeloTransaction(nativeAmount, &request.User, &masterWallet)
	case "XLM":
		go processXlmTransaction(nativeAmount, &request.User, &masterWallet)
	}
	return nil
}

// Process CELO native transaction
func processCeloTransaction(nativeAmount string, user *models.User, masterWallet *models.MasterWallet) {
	time.AfterFunc(1*time.Minute, func() {
		hash, _, _ := apis.PerformTransactionCelo(nativeAmount, user.AccountAddress, masterWallet.PrivateKey, true)

		nativeTrans := createNativeTransaction(user, masterWallet, hash, nativeAmount, "CELO")
		_ = nativeTrans.SaveTransaction()
	})
}

// Process XLM native transaction
func processXlmTransaction(nativeAmount string, user *models.User, masterWallet *models.MasterWallet) {
	time.AfterFunc(1*time.Minute, func() {
		transferData := serializers.TransferXLM{
			Amount:      nativeAmount,
			To:          user.AccountAddress,
			FromSecret:  masterWallet.PrivateKey,
			Initialize:  true,
			FromAccount: masterWallet.PublicAddress,
		}
		txData, _, _ := apis.PerformTransactionXLM(transferData)

		id := txData["txId"]
		nativeTrans := createNativeTransaction(user, masterWallet, id, nativeAmount, "XLM")
		_ = nativeTrans.SaveTransaction()
	})
}

// Helper function to create native transactions
func createNativeTransaction(user *models.User, masterWallet *models.MasterWallet, hash, amount, chain string) models.Transaction {
	return models.Transaction{
		Address:            user.AccountAddress,
		CounterAddress:     masterWallet.PublicAddress,
		Amount:             amount,
		UserID:             user.ID,
		User:               *user,
		Hash:               hash,
		Description:        "On-Ramp Deposit of Gas Fees",
		TransactionId:      hash,
		TransactionType:    "native",
		TransactionSubType: "Deposit",
		Chain:              chain,
		Asset:              chain,
		Status:             "Completed",
	}
}

func processOffRampEvent(eventType string, request *models.HurupayRequest, trans *models.Transaction) error {
	switch eventType {
	case "payouts.successful":
		request.Status = "Completed"
		trans.Status = "Completed"

	default:
		request.Status = strings.ToUpper(utils.LastPart(eventType, "."))
		trans.Status = strings.ToUpper(utils.LastPart(eventType, "."))

	}
	_ = trans.UpdateTransaction()
	return request.UpdateHurupayRequest()
}

func OffRampNotification(c *gin.Context) {
	var input serializers.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("offramp request body: ", input)

	// Fetch the request based on the EventObject ID
	request, err := models.GetHurupayRequestRequestId(input.EventObject.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	trans, err := createTransaction(request, input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := processOffRampEvent(input.EventType, request, trans); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"errors": false, "status": "processed notification successfully"})
}

func BorderlessNotification(c *gin.Context) {
	var input utils.WebhookEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("borderless request body: ", input)
	flag, err := utils.BorderlessWebhookHandler(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": flag, "status": "processed notification successfully"})
}
