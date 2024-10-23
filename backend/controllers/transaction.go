package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils"
	"backend/utils/mails"
	"backend/utils/signing"
	"backend/utils/tokens"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stellar/go/support/log"
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
	transactions, err := models.GetTransactionsByUserID(userId)
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

func GetDestinationBankAccount(c *gin.Context) {
	countryCode := c.Query("countryCode")
	root, _ := os.Getwd()

	jsonFilePath := filepath.Join(root, "templates", "bankaccount.json")

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data serializers.BankData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, bank := range data.Banks {
		if bank.CountryCode == countryCode {
			c.JSON(http.StatusOK, gin.H{
				"errors": false,
				"status": "fetch destination bank successfully",
				"data":   bank,
			})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "no available bank found for region: " + countryCode})
}

func FetchNetwork(c *gin.Context) {

	root, _ := os.Getwd()

	jsonFilePath := filepath.Join(root, "/templates", "/network.json")

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []serializers.NetworkData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "status": "fetched accepted networks for mobile money", "errors": false})
}

func GenerateReference(c *gin.Context) {
	reference := models.GenerateRequestReference()
	ref := map[string]string{
		"reference": reference,
	}
	c.JSON(200, gin.H{
		"errors": false,
		"status": "reference generated successfully",
		"data":   ref,
	})
}

func OnRampV2(c *gin.Context) {
	var input serializers.OnRamp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
	var deposit models.DepositRequest
	deposit.UserID = user.ID
	deposit.User = user
	deposit.FiatAmount = input.FiatAmount
	deposit.Currency = input.Currency
	deposit.DepositBank = input.BankName
	deposit.Ref = input.Ref
	deposit.AccountNumber = input.AccountNumber
	deposit.AccountName = input.AccountName
	deposit.AssetEquivalent = input.AssetAmount
	deposit.Status = "pending"
	deposit.ProposedAsset = input.Asset
	if err = deposit.SaveDepositRequest(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go func() {
		var adminEmails []string
		admins, err := models.FindAdmins()
		if err != nil {
			return
		}
		for _, admin := range admins {
			adminEmails = append(adminEmails, admin.Email)
		}
		fmt.Println("emails:", adminEmails)
		onRamp := serializers.AdminOnRampSerializer{
			Name:          "Admin",
			BankName:      input.BankName,
			AccountName:   input.AccountName,
			AccountNumber: input.AccountNumber,
			Amount:        input.FiatAmount,
			Currency:      input.Currency,
			Ref:           input.Ref,
		}
		_ = mails.AdminOnRampMail(adminEmails, onRamp)
	}()
	// fmt.Println(deposit)
	c.JSON(200, gin.H{
		"errors": false,
		"status": "deposit request submitted successfully",
		"data":   deposit,
	})

}

func FetchOnRampRequests(c *gin.Context) {
	ref := c.Query("ref")
	currency := c.Query("currency")
	fiatAmount := c.Query("fiat_amount")
	AccountNumber := c.Query("account_number")
	status := c.Query("status")
	countryCode := c.Query("country_code")
	cryptoAsset := c.Query("crypto_asset")
	requests, err := models.FilterDepositRequests(ref, currency, fiatAmount, AccountNumber, status, countryCode, cryptoAsset)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(requests)
	c.JSON(200, gin.H{
		"errors": false,
		"status": "requests fetched successfully",
		"data":   requests,
	})
}

func FetchOffRampRequests(c *gin.Context) {
	chain := c.Query("chain")
	hash := c.Query("hash")
	address := c.Query("address")
	AccountNumber := c.Query("account_number")
	status := c.Query("status")

	requests, err := models.FilterWithdrawalRequests(status, chain, hash, address, AccountNumber)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(requests)
	c.JSON(200, gin.H{
		"errors": false,
		"status": "requests fetched successfully",
		"data":   requests,
	})
}

func GetOnRampRequest(c *gin.Context) {
	id := c.Param("id")
	intId, _ := strconv.Atoi(id)
	deposit, err := models.GetDepositRequest(intId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"errors": false,
		"status": "request fetched successfully",
		"data":   deposit,
	})
}

func GetOffRampRequest(c *gin.Context) {
	id := c.Param("id")
	intId, _ := strconv.Atoi(id)
	withdrawal, err := models.GetWithdrawalRequest(intId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"errors": false,
		"status": "request fetched successfully",
		"data":   withdrawal,
	})
}

func VerifyOnRamp(c *gin.Context) {
	id := c.Param("id")
	intId, _ := strconv.Atoi(id)
	deposit, err := models.GetDepositRequest(intId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "message": "error fetching request"})
		return
	}
	var input serializers.OnRampAction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "message": "in"})
		return
	}
	user, err := models.GetUserByID(deposit.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	masterWallet, err := models.FetchMasterWallet(user.CryptoCurrency)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	nativeAmount, err := utils.PerformDepositofNativeCalculation(deposit.AssetEquivalent, "USD", user.CryptoCurrency)
	// fmt.Println(nativeAmount)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	switch input.Action {
	case "Approve":
		deposit.Status = "Approved"
		var transaction models.Transaction
		switch user.CryptoCurrency {
		case "CELO":
			txHash, code, err := apis.PerformTransactionCelo(deposit.AssetEquivalent, user.AccountAddress, masterWallet.PrivateKey, false)
			if err != nil {
				c.JSON(code, gin.H{"error": err.Error(), "message": "transaction failed"})
				return
			}
			transaction.Hash = txHash
			transaction.Chain = "CELO"
			go func() {
				time.AfterFunc(1*time.Minute, func() {
					hash, _, _ := apis.PerformTransactionCelo(nativeAmount, user.AccountAddress, masterWallet.PrivateKey, true)

					var nativeTrans models.Transaction
					nativeTrans.Address = user.AccountAddress
					nativeTrans.CounterAddress = masterWallet.PublicAddress
					nativeTrans.Amount = nativeAmount
					nativeTrans.UserID = user.ID
					nativeTrans.User = user
					nativeTrans.Hash = hash
					nativeTrans.Description = "On-Ramp Deposit of Gas Fees"
					nativeTrans.TransactionId = hash
					nativeTrans.TransactionType = "native"
					nativeTrans.TransactionSubType = "Deposit"
					nativeTrans.Chain = "CELO"
					nativeTrans.Asset = "CELO"
					nativeTrans.Status = "Completed"
					_ = nativeTrans.SaveTransaction()

				})
			}()

		case "XLM":
			transferData := serializers.TransferXLM{
				Amount:        deposit.AssetEquivalent,
				To:            user.AccountAddress,
				FromSecret:    masterWallet.PrivateKey,
				Initialize:    true,
				Token:         deposit.ProposedAsset,
				IssuerAccount: masterWallet.PublicAddress,
				FromAccount:   masterWallet.PublicAddress,
			}
			txData, code, err := apis.PerformTransactionXLM(transferData)
			if err != nil {
				c.JSON(code, gin.H{"error": err.Error(), "message": "transaction failed"})
				return
			}
			hash := txData["txId"]
			transaction.Hash = hash
			transaction.Chain = "XLM"

			go func() {
				time.AfterFunc(1*time.Minute, func() {
					transferData := serializers.TransferXLM{
						Amount:      nativeAmount,
						To:          user.AccountAddress,
						FromSecret:  masterWallet.PrivateKey,
						Initialize:  true,
						FromAccount: masterWallet.PublicAddress,
					}
					txData, _, _ := apis.PerformTransactionXLM(transferData)

					var xlmTrans models.Transaction
					id := txData["txId"]
					xlmTrans.Hash = id
					xlmTrans.Chain = "XLM"
					xlmTrans.Amount = nativeAmount
					xlmTrans.UserID = user.ID
					xlmTrans.User = user
					xlmTrans.TransactionId = id
					xlmTrans.TransactionType = "native"
					xlmTrans.TransactionSubType = "Deposit"
					xlmTrans.Asset = "XLM"
					xlmTrans.Description = "On-Ramp Deposit of Gas Fees"
					xlmTrans.Status = "Completed"
					_ = xlmTrans.SaveTransaction()

				})
			}()

		}

		transaction.Amount = deposit.AssetEquivalent
		transaction.Status = "Completed"

		transaction.Asset = deposit.ProposedAsset
		transaction.Address = user.AccountAddress
		transaction.CounterAddress = masterWallet.PublicAddress
		transaction.TransactionSubType = "Deposit"
		transaction.TransactionType = "Fungible Token"
		transaction.Description = "On-Ramp Deposit"
		transaction.User = user
		transaction.UserID = user.ID
		if err := transaction.SaveTransaction(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	case "Reject":
		deposit.Status = "Rejected"
	}
	_ = deposit.UpdateDepositRequest()
	c.JSON(200, gin.H{
		"errors": false,
		"status": "request verified successfully",
		"data":   deposit,
	})
}

func OffRampV2(c *gin.Context) {
	var input serializers.OffRamp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	withdrawal := models.WithdrawalRequest{
		UserID:         user.ID,
		User:           user,
		CryptoAmount:   input.CryptoAmount,
		BankName:       input.BankName,
		AccountNumber:  input.AccountNumber,
		AccountName:    input.AccountName,
		Address:        user.AccountAddress,
		Status:         "pending",
		Asset:          input.Asset,
		Chain:          input.Chain,
		EquivalentFiat: input.FiatEquivalent,
		FiatCurrency:   input.CurrencyCode,
	}

	masterWallet, err := models.FetchMasterWallet(input.Chain)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	trans := models.Transaction{}
	if err := processTransaction(&trans, &withdrawal, input, user, masterWallet); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error(), "message": "transaction failed"})
		return
	}
	trans.UserID = user.ID
	trans.User = user
	trans.Amount = withdrawal.CryptoAmount
	trans.TransactionType = "Fungible Token"
	trans.TransactionSubType = "Withdrawal"
	if err := trans.SaveTransaction(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err = withdrawal.SaveWithdrawalRequest(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go notifyAdmins(input, trans.Hash)

	c.JSON(200, gin.H{
		"errors": false,
		"status": "withdrawal request submitted successfully",
		"data":   withdrawal,
	})
}

func processTransaction(trans *models.Transaction, withdrawal *models.WithdrawalRequest, input serializers.OffRamp, user models.User, masterWallet models.MasterWallet) *apiError {
	switch input.Chain {
	case "CELO":
		hash, code, err := apis.PerformTransactionCelo(withdrawal.CryptoAmount, masterWallet.PublicAddress, user.PrivateKey, false)
		if err != nil {
			return &apiError{code, err}
		}
		trans.Hash = hash
		trans.Chain = "CELO"
		trans.Asset = "CUSD"
		withdrawal.Status = "Awaiting Payment"
		withdrawal.Hash = hash
	case "XLM":
		transferData := serializers.TransferXLM{
			Amount:        withdrawal.CryptoAmount,
			To:            masterWallet.PublicAddress,
			FromSecret:    user.PrivateKey,
			Initialize:    false,
			Token:         input.Asset,
			IssuerAccount: user.AccountAddress,
			FromAccount:   user.AccountAddress,
		}
		txData, code, err := apis.PerformTransactionXLM(transferData)
		if err != nil {
			return &apiError{code, err}
		}
		hash := txData["txId"]
		trans.Hash = hash
		trans.Chain = "XLM"
		trans.Asset = input.Asset
		withdrawal.Status = "Awaiting Payment"
		withdrawal.Hash = hash
	}
	return nil
}

func notifyAdmins(input serializers.OffRamp, hash string) {
	admins, err := models.FindAdmins()
	if err != nil {
		return
	}

	var adminEmails []string
	for _, admin := range admins {
		adminEmails = append(adminEmails, admin.Email)
	}

	offRamp := serializers.AdminOffRampSerializer{
		Name:          input.AccountName,
		BankName:      input.BankName,
		AccountNumber: input.AccountNumber,
		Amount:        input.FiatEquivalent,
		Currency:      input.CurrencyCode,
		Ref:           hash,
	}
	_ = mails.AdminOffRampMail(adminEmails, offRamp)
}

type apiError struct {
	StatusCode int
	Err        error
}

func (e *apiError) Error() string {
	return e.Err.Error()
}

func VerifyOffRamp(c *gin.Context) {
	id := c.Param("id")
	intId, _ := strconv.Atoi(id)
	withdrawal, err := models.GetWithdrawalRequest(intId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var input serializers.OffRampAction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	withdrawal.Status = "Completed"
	withdrawal.BankRef = input.BankRef
	if err := withdrawal.UpdateWithdrawalRequest(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserByID(withdrawal.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	go func() {
		floatAmount, err := strconv.ParseFloat(withdrawal.EquivalentFiat, 64)
		if err != nil {
			panic(err)
		}
		amount := utils.FormatAmountWithCommas(floatAmount)
		data := serializers.UserOffRampMail{
			Name:          fmt.Sprintf("%s %s", user.LastName, user.FirstName),
			Amount:        amount,
			Currency:      withdrawal.FiatCurrency,
			Ref:           withdrawal.BankRef,
			BankName:      withdrawal.BankName,
			AccountNumber: withdrawal.AccountNumber,
			AccountName:   withdrawal.AccountName,
		}
		emails := []string{user.Email}
		_ = mails.UserOffRampMail(emails, data)
	}()
	_ = withdrawal.UpdateWithdrawalRequest()
	c.JSON(200, gin.H{
		"errors": false,
		"status": "request verified successfully",
	})

}

func GetExchangeRate(c *gin.Context) {
	fiatCurrency := c.Query("fiat_currency")
	asset := c.Query("asset")
	rate := make(chan string)
	err := make(chan error)
	go apis.GetExchangeRate(fiatCurrency, asset, rate, err)

	select {
	case exchangeRate := <-rate:
		c.JSON(200, gin.H{
			"errors": false,
			"status": "exchange rate fetched successfully",
			"data":   exchangeRate,
		})
	case error := <-err:
		c.JSON(400, gin.H{
			"errors": true,
			"status": error.Error(),
		})
	case <-time.After(10 * time.Second):
		c.JSON(408, gin.H{
			"errors": true,
			"status": "Request timed out",
		})
	}
}

func MobileMoneyAmountToReceive(c *gin.Context) {
	amount := c.Query("amount")
	currency := c.Query("currency")
	asset := c.Query("cryptoAsset")
	transType := c.Query("type")
	rateChan := make(chan string)
	errChan := make(chan error)
	go apis.GetMobileMoneyExhangeRate(currency, rateChan, errChan)
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
	case <-time.After(15 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
		return
	}

}

func MobileMoneyOnRamp(c *gin.Context) {
	var input serializers.Payment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.Transfer.DigitalAsset == "CUSD" {
		input.Transfer.DigitalAsset = "cUSD"
	}

	resp, err := apis.OnRampMobileMoney(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var request models.HurupayRequest
	request.Amount = strconv.FormatInt(int64(input.Collection.Amount), 10)
	request.UserId = int32(user.ID)
	request.AccountNumber = input.Collection.PhoneNumber
	request.Status = "Pending"
	request.RequestType = models.OnRamp
	request.MobileNumber = input.Collection.PhoneNumber
	request.Token = input.Transfer.DigitalAsset
	request.MobileNetwork = input.Collection.Network
	request.CryptoChain = input.Transfer.DigitalNetwork
	request.RequestId = resp.Data.CollectionRequestID
	request.CountryCurrency = input.Collection.CountryCode
	if err := request.SaveHurupayRequest(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"errors": false,
		"status": "mobile money on ramp initiated",
		"data":   resp,
	})

}

func MobileMoneyOffRamp(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}

	var input serializers.MobileOffRamp
	if err := c.ShouldBindJSON(&input); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateOffRampRequest(input, user); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	data := createTransactionRequest(input)
	resp, err := apis.OffRampMobileMoney(data)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := storeHurupayRequest(input, user, resp.Data.PayoutRequestID); err != nil {
		log.Error(err)
	}

	go handleOffRampTransaction(input, user, resp)

	c.JSON(http.StatusOK, gin.H{
		"errors": false,
		"status": "mobile money off ramp initiated",
	})
}

// Helper to send error responses
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// Validate the input and user's account address
func validateOffRampRequest(input serializers.MobileOffRamp, user models.User) error {
	if input.Token == "CUSD" {
		input.Token = "cUSD"
	}
	if user.AccountAddress != input.SendingAddress {
		return errors.New("invalid account address")
	}
	return nil
}

// Create the transaction request data
func createTransactionRequest(input serializers.MobileOffRamp) serializers.TransactionRequest {
	return serializers.TransactionRequest{
		SendingAddress: input.SendingAddress,
		AmountSending:  input.AmountSending,
		Network:        input.Network,
		Token:          input.Token,
	}
}

// Handle the transaction in a goroutine
func handleOffRampTransaction(input serializers.MobileOffRamp, user models.User, resp apis.PayoutResponse) {
	var hash string
	var err error

	switch input.Network {
	case "CELO":
		hash, err = executeCeloTransaction(input, user, resp.Data.EscrowAddress)
	case "XLM":
		hash, err = executeXlmTransaction(input, user, resp.Data.EscrowAddress)
	}

	if err != nil {
		log.Error(err)
	}
	fmt.Println("hash: ", hash)

	transaction := prepareTransactionDetails(input, resp, hash)

	output, err := apis.OffRampMobileFinalize(transaction)
	if err != nil || output.Data.ResultCode != 0 {
		log.Error("hurrupay failed to finalize", err)
		return
	}

}

// Execute the CELO transaction
func executeCeloTransaction(input serializers.MobileOffRamp, user models.User, escrowAddress string) (string, error) {
	hash, _, err := apis.PerformTransactionCelo(input.AmountSending, escrowAddress, user.PrivateKey, false)
	return hash, err
}

// Execute the XLM transaction
func executeXlmTransaction(input serializers.MobileOffRamp, user models.User, escrowAddress string) (string, error) {
	transferData := serializers.TransferXLM{
		Amount:        input.AmountSending,
		To:            escrowAddress,
		FromSecret:    user.PrivateKey,
		Initialize:    false,
		Token:         "USDC",
		IssuerAccount: user.AccountAddress,
		FromAccount:   user.AccountAddress,
	}
	txData, _, err := apis.PerformTransactionXLM(transferData)
	if err != nil {
		return "", err
	}
	txID := txData["txId"]
	return txID, nil
}

// Prepare the transaction details
func prepareTransactionDetails(input serializers.MobileOffRamp, resp apis.PayoutResponse, hash string) serializers.TransactionDetails {
	return serializers.TransactionDetails{
		Collection: struct {
			TransactionHash string `json:"transactionHash"`
			PayoutRequestID string `json:"payoutRequestId"`
			Network         string `json:"network"`
			Token           string `json:"token"`
		}{
			TransactionHash: hash,
			PayoutRequestID: resp.Data.PayoutRequestID,
			Network:         input.Network,
			Token:           input.Token,
		},
		Transfer: struct {
			CustomerName string `json:"customerName"`
			PhoneNumber  string `json:"phoneNumber"`
			CountryCode  string `json:"countryCode"`
			Network      string `json:"network"`
		}{
			CustomerName: input.CustomerName,
			PhoneNumber:  input.PhoneNumber,
			CountryCode:  input.CountryCode,
			Network:      input.MobileProvider,
		},
	}
}

// Store the Hurupay request to the database
func storeHurupayRequest(input serializers.MobileOffRamp, user models.User, requestId string) error {
	request := models.HurupayRequest{
		Amount:          input.AmountSending,
		UserId:          int32(user.ID),
		AccountNumber:   input.PhoneNumber,
		Status:          "Pending",
		RequestType:     models.OffRamp,
		MobileNumber:    input.PhoneNumber,
		Token:           input.Token,
		MobileNetwork:   input.MobileProvider,
		CryptoChain:     input.Network,
		RequestId:       requestId,
		CountryCurrency: input.CountryCode,
	}
	return request.SaveHurupayRequest()
}
