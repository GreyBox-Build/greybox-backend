package controllers

import (
	//"backend/apis"
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/mails"
	"backend/utils/tokens"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func CreateAccountV2(c *gin.Context) {

	var input serializers.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	var user models.User
	user.FirstName = input.FirstName
	user.Email = input.Email
	user.Password = input.Password
	user.LastName = input.LastName
	user.Country = input.Country
	user.LastName = input.LastName
	user.Currency = input.Currency
	user.CountryCode = input.CountryCode
	if err := user.BeforeSaveDetail(); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.CryptoCurrency = input.Chain
	masterWallet, err := models.FetchMasterWallet(input.Chain)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.SignatureId = masterWallet.SignatureId
	user.Index = masterWallet.CurrentIndex + 1
	masterWallet.TotalAddressesGenerated += 1

	privData, err := apis.GetPrivateKeyManagedWallet(masterWallet.SignatureId, user.Index)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.PrivateKey = privData.PrivateKey
	addressData, err := apis.GetManagedWalletAddress(masterWallet.SignatureId, user.Index)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.AccountAddress = addressData.Address
	switch input.Chain {
	case serializers.Chains.Celo:
		user.CryptoCurrency = serializers.Chains.Celo
		walletData, err := apis.GetManagedWallet(masterWallet.SignatureId)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		user.Xpub = walletData.Xpub
		user.Mnemonic = walletData.Mnemonic

	case serializers.Chains.Stellar:
		user.CryptoCurrency = serializers.Chains.Stellar
	}

	if err := user.SaveUser(); err != nil {
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"message": "creating user falied",
		})
		return
	}
	masterWallet.TotalAddressActivated += 1
	if err := masterWallet.UpdateMasterWallet(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "created account successfuly",
	})
}

func FetchAuthenticatedUserToken(c *gin.Context) {
	var input serializers.LoginSerializer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	email := input.Email
	password := input.Password
	token, err := models.LoginCheck(email, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	data := map[string]string{
		"access_token": token,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "fetched access token",
		"data":   data,
		"errors": false,
	})

}

func ForgetPassword(c *gin.Context) {
	// Parse the request JSON
	var requestData struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email exists in the database
	user, ok := models.FindUserByEmail(requestData.Email)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate a unique token
	token, err := models.GenerateRecoveryToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	receiver := []string{
		user.Email,
	}

	if err := mails.SendForgetPasswordMail(receiver, user.FirstName, token.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func ResetPassword(c *gin.Context) {
	// Parse the request JSON
	var requestData struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := models.CheckTokenValid(requestData.Token)
	// Check if the token exists and is valid
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update the user's password in the database
	user, err := models.GetUserByID(token.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Password = requestData.Password
	if err := user.BeforeSaveDetail(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	user.UpdateUser()

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func GetAuthenticatedUser(c *gin.Context) {
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
	var balance float32
	switch user.CryptoCurrency {
	case serializers.Chains.Celo:
		balance, err = apis.FetchAccountBalanceCUSD(user.AccountAddress)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

	case serializers.Chains.Stellar:
		balance, err = apis.FetchAccountBalanceXLM(user.AccountAddress)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	data := map[string]float32{
		"balance": balance,
	}
	authData := map[string]interface{}{
		"personal_details": user,
		"wallet_details":   data,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "fetch authenticated user details",
		"errors": false,
		"data":   authData,
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

func CreateMasterWallet(c *gin.Context) {
	var input serializers.MasterWalletForm
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var wallet models.MasterWallet
	signature, err := apis.StorePrivateKeyManagedWallet(input.Asset)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	wallet.SignatureId = signature.SignatureId
	privKey, err := apis.GetPrivateKeyManagedWallet(signature.SignatureId, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	wallet.PrivateKey = privKey.PrivateKey
	switch input.Asset {
	case serializers.Chains.Celo:
		walletData, err := apis.GetManagedWallet(wallet.SignatureId)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		wallet.XpublicAddress = walletData.Xpub
		wallet.Mnemonic = walletData.Mnemonic

		wallet.WalletChain = serializers.Chains.Celo

	case serializers.Chains.Stellar:
		wallet.WalletChain = serializers.Chains.Stellar
	}
	walletAddress, err := apis.GetManagedWalletAddress(wallet.SignatureId, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	wallet.PublicAddress = walletAddress.Address
	if err := wallet.CreateMasterWallet(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	data := map[string]interface{}{
		"account_address": wallet.PublicAddress,
		"asset":           wallet.WalletChain,
	}
	c.JSON(200, gin.H{
		"status": "created master wallet",
		"data":   data,
		"errors": false,
	})

}

func FetchMasterWallet(c *gin.Context) {
	chain := c.Query("chain")
	wallet, err := models.FetchMasterWallet(chain)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	data := map[string]interface{}{
		"account_address": wallet.PublicAddress,
		"asset":           wallet.WalletChain,
	}
	c.JSON(200, gin.H{
		"status": "fetched master wallet",
		"data":   data,
		"errors": false,
	})

}
