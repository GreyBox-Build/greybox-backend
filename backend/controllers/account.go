package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/mails"
	"backend/utils/tokens"
	"net/http"
	"os"
	"strings"

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

	meumnic, xpub, err := apis.GenerateCelloWallet()
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "generating cello wallet failed",
		})
		return
	}
	address, err := apis.GenerateCelloAddress(xpub)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "generating address failed",
		})
		return
	}
	user.AccountAddress = address

	user.Xpub = xpub
	user.Mnemonic = meumnic
	privURL := "https://api.tatum.io/v3/celo/wallet/priv"

	privData := serializers.PrivGeneration{
		Index:    1,
		Mnemonic: meumnic,
	}
	key := os.Getenv("TATUM_API_KEY_TEST")
	privKey, err := apis.GeneratePrivateKey(privURL, key, privData)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "generating private key failed",
		})
		return
	}

	user.PrivateKey = privKey

	if err := user.SaveUser(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "creating user falied",
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
	balance, err := apis.FetchWalletBalance(user.AccountAddress, strings.ToLower(user.CryptoCurrency), 10)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
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
