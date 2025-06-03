package controllers

import (
	"backend/apis"
	"backend/apis/borderless"
	"backend/models"
	"backend/serializers"
	"backend/state"
	"backend/utils"
	"backend/utils/tokens"
	"fmt"
	"strings"

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
	borderless := borderless.NewBorderless()
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
	paymentInstruction := borderless.NewPayment(
		bank.Country, input.Currency, fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		"Wire", input.AccountHolderName, input.AccountNumber,
		input.AccountType, bank.Name, *bank.Street, *bank.City, bank.Country, *bank.ZipCode, *bank.SwiftCode,
		input.AccountNumber, *bank.SwiftCode, *bank.Street, *bank.State)

	borderlessHandler := borderless.NewBorderless()
	paymentInstructionResponse, err := borderlessHandler.MakePaymentInstruction(paymentInstruction)
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
	accountID := state.AppConfig.BorderlessAccountId
	withdraw := borderless.NewWithdrawalRequest(
		input.Currency, bank.Country, "USDC_POLYGON", input.Amount,
		accountID, input.PaymentPurpose, paymentInstructionResponse.ID)
	res, err := borderlessHandler.MakeWithdrawal(withdraw)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	borderlessRequest.TxId = res.ID
	if err := models.CreateBorderlessRequest(&borderlessRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transactionInstruction, err := borderlessHandler.GetTransaction(res.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": transactionInstruction, "status": "success", "errors": false})

}

func BorderlessMobileMoneyOnRamp(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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

	// make sure country in request is a valid alpha-2 code
	validCodes := utils.CreateValidCountryCodes()
	country := strings.ToUpper(input.Country)

	if !validCodes[country] {
		c.JSON(400, gin.H{"error": "Invalid country code"})
		return
	}

	borderless := borderless.NewBorderless()
	availableCountries, err := borderless.GetAvailableCountries("deposits")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// make sure response contains the country in the request
	for _, v := range availableCountries {
		if v == strings.ToUpper(country) {
			break
		}
		c.JSON(400, gin.H{"error": "Country not supported"})
		return
	}

	// find mobile money deposit option
	depositOption, err := borderless.GetDepositOrWithdrawalOption("deposits", input.Country, input.Fiat, input.Asset)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if depositOption == nil {
		c.JSON(400, gin.H{"error": "Deposit option not found"})
		return
	}

	makeDepositResponse, err := borderless.MobileMoneyDeposit(
		input.AccountId, input.Fiat, input.Country, input.Asset, input.Amount, depositOption.Method)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	borderlessRequest := models.BorderlessRequest{}
	borderlessRequest.FiatAmount = input.Amount
	borderlessRequest.Asset = input.Asset
	borderlessRequest.Country = input.Country
	borderlessRequest.UserId = userId
	borderlessRequest.User = user
	borderlessRequest.Status = "Pending"
	borderlessRequest.AccountId = makeDepositResponse.Destination.AccountID
	borderlessRequest.TxId = makeDepositResponse.ID
	borderlessRequest.FeeAmount = makeDepositResponse.FeeAmount

	if err := models.CreateBorderlessRequest(&borderlessRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transactionInstruction, err := borderless.GetTransaction(makeDepositResponse.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": transactionInstruction, "status": "success", "errors": false})
}

func BorderlessMobileMoneyOffRamp(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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

	// if asset is not provided, throw error
	if input.Asset == "" {
		c.JSON(400, gin.H{"error": "Asset is required"})
		return
	}

	// if master wallet is not provided, throw error
	if input.MasterWallet == "" {
		c.JSON(400, gin.H{"error": "Master wallet is required"})
		return
	}

	bank, err := models.GetBankData(int(input.BankId))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	borderlessHandler := borderless.NewBorderless()

	availableCountries, err := borderlessHandler.GetAvailableCountries("withdrawals")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// make sure mobile money is available for the bank's country
	for _, v := range availableCountries {
		if v == strings.ToUpper(bank.Country) {
			break
		}
		c.JSON(400, gin.H{"error": "Country not supported"})
		return
	}

	// find mobile money withdrawal option
	withdrawalOption, err := borderlessHandler.GetDepositOrWithdrawalOption("withdrawals", bank.Country, input.Currency, input.Asset)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if withdrawalOption == nil {
		c.JSON(400, gin.H{"error": "Deposit option not found"})
		return
	}

	// create payment instruction
	paymentInstruction := borderless.NewPayment(
		bank.Country, input.Currency, fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		withdrawalOption.Method, input.AccountHolderName, input.AccountNumber, input.AccountType,
		bank.Name, *bank.Street, *bank.City, bank.Country, *bank.ZipCode,
		*bank.SwiftCode, input.AccountNumber, *bank.SwiftCode, *bank.Street, *bank.State)

	paymentInstructionResponse, err := borderlessHandler.MakePaymentInstruction(paymentInstruction)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	masterWallet, err := models.FetchMasterWallet(input.MasterWallet)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var trans models.Transaction
	trans.UserID = userId
	trans.User = user
	trans.Amount = input.Amount
	trans.Status = "pending"
	trans.Chain = input.MasterWallet
	trans.Address = user.AccountAddress
	trans.TransactionSubType = "Withdrawal"
	trans.TransactionType = "fungible"

	var borderlessRequest models.BorderlessRequest
	borderlessRequest.UserId = userId
	borderlessRequest.User = user
	borderlessRequest.Status = "Pending"
	borderlessRequest.FiatAmount = input.Amount
	borderlessRequest.Asset = input.Asset
	borderlessRequest.PaymentInstructionId = &paymentInstructionResponse.ID

	currency, err := apis.ParseCurrencyType(input.Asset)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tatumInstance := apis.NewTatumPolygon()
	hashResponse, err := tatumInstance.PerformTransaction(
		masterWallet.PublicAddress, input.Amount,
		user.PrivateKey, currency)
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

	accountID := state.AppConfig.BorderlessAccountId
	if input.AccountId != "" {
		accountID = input.AccountId
	}

	withdraw := borderless.NewWithdrawalRequest(
		input.Currency, bank.Country, input.Asset,
		input.Amount, accountID, input.PaymentPurpose, paymentInstructionResponse.ID)
	res, err := borderlessHandler.MakeWithdrawal(withdraw)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	borderlessRequest.TxId = res.ID
	if err := models.CreateBorderlessRequest(&borderlessRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transactionInstruction, err := borderlessHandler.GetTransaction(res.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": transactionInstruction, "status": "success", "errors": false})

}
