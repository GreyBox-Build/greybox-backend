package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/mails"
	"backend/utils/tokens"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {

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
	if err := user.BeforeSave(); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.AccountID = models.GenerateAccountId()
	user.AccountNumber = models.GenerateAccountNumber()
	user.AccountCode = models.GenerateAccountCode("GBX")

	apiURL := "https://api.tatum.io/v3/ledger/account"
	key := os.Getenv("TATUM_API_KEY_TEST")
	customer := map[string]string{
		"externalId":         user.AccountID,
		"accountingCurrency": user.Currency,
		"customerCountry":    user.CountryCode,
		"providerCountry":    "GH",
	}
	xpub, err := models.GetLatestXlmPublic()
	if err != nil {
		data, err := apis.GenerateStellarAddress()
		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "generating address failed",
			})
			return
		}
		xpub.Xpub = data["address"]
		xpub.Secret = data["secret"]
		if err := xpub.Save(); err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "failed to create xpub key",
			})
		}
	}

	virtual := serializers.VirtualAccount{
		Xpub:               xpub.Xpub,
		Currency:           "XLM",
		Customer:           customer,
		Compliant:          false,
		AccountCode:        user.AccountCode,
		AccountNumber:      user.AccountNumber,
		AccountingCurrency: user.Currency,
	}
	user.Xpub = xpub.Xpub
	user.PrivateKey = xpub.Secret
	id, err := apis.CreateVirtualAccount(apiURL, key, virtual)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "creating virtual account failed",
		})
		return
	}
	user.CustomerId = id
	address, err := apis.CreateDepositWalletXLM(user.CustomerId)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "creating deposit wallet failed",
		})
		return
	}
	user.AccountAddress = address

	if err := user.SaveUser(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "creating user falied",
		})
		return
	}
	if err := apis.ActivateVirtualAccount(user.AccountID, key); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "activating virtual account failed",
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
	if err := user.BeforeSave(); err != nil {
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
	apiKey := os.Getenv("TATUM_API_KEY_TEST")
	data, err := apis.FetchAccountBalance(user.CustomerId, apiKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
