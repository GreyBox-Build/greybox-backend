package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/tokens"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func BorderLessOnramp(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var input serializers.BorderlessOnramp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	borderless := apis.NewBorderless()
	response, err := borderless.MakeDeposit(
		input.Amount, input.Asset, input.Country, input.Fiat)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	borderlessRequest := models.BorderlessRequest{}
	borderlessRequest.FiatAmount = input.Amount
	borderlessRequest.Asset = input.Asset     // change this for prod use
	borderlessRequest.Country = input.Country //change this for prod use
	borderlessRequest.UserId = userId
	borderlessRequest.User = user
	borderlessRequest.Status = "Pending"
	borderlessRequest.AccountId = response.Destination.AccountID
	borderlessRequest.TxId = response.ID
	borderlessRequest.FeeAmount = response.FeeAmount
	if err := models.CreateBorderlessRequest(&borderlessRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	transactionInstruction, err := borderless.GetTransaction(response.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": transactionInstruction, "status": "success", "errors": false})
}

func FilterBank(c *gin.Context) {
	country := c.Query("country")

	banks, err := models.FilterBanks(country)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": banks, "status": "success", "errors": false})
}

func BorderLessOffRamp(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var input serializers.MakeWithdrawalBorderless
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	bank, err := models.GetBankData(int(input.BankId))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	polygon := apis.NewTatumPolygon()
	borderless := apis.NewBorderless()
	paymentInstruction := apis.NewPayment(
		bank.Country, input.Currency, fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		"Wire", input.AccountHolderName, input.AccountNumber,
		input.AccountType, bank.Name, *bank.Street, *bank.City, bank.Country, *bank.ZipCode, *bank.SwiftCode,
		input.AccountNumber, *bank.SwiftCode, *bank.Street, *bank.State)
	paymentInstructionResponse, err := borderless.MakePaymentInstruction(paymentInstruction)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	masterWallet, err := models.FetchMasterWallet("MATIC")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	var trans models.Transaction
	trans.UserID = userId
	trans.User = user
	trans.Amount = input.Amount
	trans.Status = "pending"
	trans.Chain = "MATIC"
	trans.Address = user.AccountAddress
	trans.TransactionSubType = "Withdrawal"
	trans.TransactionType = "fungible"
	var borderlessRequest models.BorderlessRequest
	borderlessRequest.UserId = userId
	borderlessRequest.User = user
	borderlessRequest.Status = "Pending"
	borderlessRequest.FiatAmount = input.Amount
	borderlessRequest.Asset = "USDC_MATIC"
	borderlessRequest.PaymentInstructionId = &paymentInstructionResponse.ID
	currency, err := apis.ParseCurrencyType("USDC_MATIC")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	hashResponse, err := polygon.PerformTransaction(
		masterWallet.PublicAddress, input.Amount, user.PrivateKey, currency)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	trans.Hash = hashResponse.TxId
	trans.TransactionId = hashResponse.TxId
	if err := trans.SaveTransaction(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	accountID := os.Getenv("BORDERLESS_ACCOUNT_ID")
	withdraw := apis.NewWithdrawalRequest(
		input.Currency, bank.Country, "USDC_POLYGON", input.Amount,
		accountID, input.PaymentPurpose, paymentInstructionResponse.ID)
	res, err := borderless.MakeWithdrawal(withdraw)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	borderlessRequest.TxId = res.ID
	if err := models.CreateBorderlessRequest(&borderlessRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transactionInstruction, err := borderless.GetTransaction(res.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": transactionInstruction, "status": "success", "errors": false})

}
