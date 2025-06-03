package controllers

import (
	"backend/apis"
	"backend/apis/borderless"
	"backend/models"
	"backend/serializers"
	"backend/state"
	"backend/utils"
	"backend/utils/mails"
	"backend/utils/tokens"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateAccountV2(c *gin.Context) {
	var input serializers.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err, "invalid request payload")
		return
	}

	// check if user already exists
	_, exists := models.FindUserByEmail(input.Email)
	if exists {
		utils.BadRequest(c, errors.New("DUPLICATE USER"), "User already exists")
		return
	}

	user := models.User{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       input.Password,
		Country:        input.Country,
		Currency:       input.Currency,
		CountryCode:    input.CountryCode,
		CryptoCurrency: input.Chain,
	}

	switch strings.ToUpper(input.Chain) {
	case serializers.Chains.Celo:
		if err := apis.SetupCeloAccount(&user); err != nil {
			utils.BadRequest(c, err, "celo wallet setup failed")
			return
		}
	case serializers.Chains.Stellar:
		if err := apis.SetupStellarAccount(&user); err != nil {
			utils.BadRequest(c, err, "stellar account setup failed")
			return
		}
	case serializers.Chains.Polygon:
		if err := apis.SetupPolygonAccount(&user); err != nil {
			utils.BadRequest(c, err, "polygon account setup failed")
			return
		}
	default:
		utils.BadRequest(c, fmt.Errorf("unsupported chain: %s", input.Chain), "unsupported chain type")
		return
	}

	if err := user.SaveUser(); err != nil {
		utils.BadRequest(c, err, "creating user failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "created account successfully"})
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
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Password reset failed",
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

	balanceChan := make(chan float32)
	errorChan := make(chan error)
	switch user.CryptoCurrency {
	case serializers.Chains.Celo:
		go func() {
			balance, err := apis.FetchAccountBalanceCUSD(user.AccountAddress)
			if err != nil {
				errorChan <- err
				return
			}
			balanceChan <- balance
		}()
	case serializers.Chains.Stellar:
		go func() {
			balance, err := apis.FetchAccountBalanceXLM(user.AccountAddress)
			if err != nil {
				errorChan <- err
				return
			}
			balanceChan <- balance
		}()
	case serializers.Chains.Polygon:
		go func() {
			balance, err := apis.FetchWalletBalance(user.AccountAddress, "polygon", 10)
			if err != nil {
				errorChan <- err
				return
			}
			balanceChan <- balance
		}()
	}

	// Wait for balance or error
	var balance float32
	select {
	case balance = <-balanceChan:
		// Balance successfully retrieved
	case err = <-errorChan:
		// Error occurred while fetching balance
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.PreviousBalance = balance
	user.UpdateUser()

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

func GenerateMasterWallet(c *gin.Context) {
	var input serializers.MasterWalletForm
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	var masterWallet models.MasterWallet
	switch input.Asset {
	case serializers.Chains.Celo:
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
		masterWallet.PublicAddress = address
		masterWallet.XpublicAddress = xpub
		masterWallet.Mnemonic = meumnic
		privData := serializers.PrivGeneration{
			Index:    1,
			Mnemonic: meumnic,
		}
		privKey, err := apis.GeneratePrivateKey(privData)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "generating private key failed",
			})
			return
		}
		masterWallet.PrivateKey = privKey

	case serializers.Chains.Stellar:
		data, code, err := apis.GenerateXlmAccount()
		if err != nil {
			c.JSON(code, gin.H{
				"error":   err.Error(),
				"message": "generating address failed",
			})
			return
		}
		address, secret := data["address"], data["secret"]
		masterWallet.PublicAddress = address
		masterWallet.PrivateKey = secret
	case serializers.Chains.Polygon:
		polygon := apis.NewTatumPolygon()
		walletResponse, err := polygon.CreateWallet()
		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "generating polygon wallet failed",
			})
			return
		}
		xPub, mnemenic := walletResponse.Xpub, walletResponse.Mnemonic
		masterWallet.XpublicAddress = xPub
		masterWallet.Mnemonic = mnemenic
		privResponse, err := polygon.GeneratePrivateKey(mnemenic, 0)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "generating polygon private key failed",
			})
			return
		}
		masterWallet.PrivateKey = privResponse.Key
		addressResponse, err := polygon.GenerateAddress(xPub, 0)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "generating polygon address failed",
			})
			return
		}
		masterWallet.PublicAddress = addressResponse.Address
	}
	masterWallet.WalletChain = input.Asset
	if err := masterWallet.CreateMasterWallet(); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "creating master wallet failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "created master wallet succesfully",
		"errors": false,
		"data":   masterWallet,
	})
}

func GetMasterWallet(c *gin.Context) {
	input := c.Query("asset")
	masterWallet, err := models.FetchMasterWallet(input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "fetched master wallet succesfully",
		"errors": false,
		"data":   masterWallet,
	})
}

func GetMasterWallets(c *gin.Context) {
	masterWallet, err := models.FetchMasterWallets()
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "fetched master wallets succesfully",
		"errors": false,
		"data":   masterWallet,
	})
}

func MakeAdmin(c *gin.Context) {
	var input serializers.AdminForm
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	key := state.AppConfig.AdminKey
	if key != input.Key {
		c.JSON(400, gin.H{
			"error": "invalid admin key",
		})
		return
	}
	user, ok := models.FindUserByEmail(input.UserEmail)
	if !ok {
		c.JSON(400, gin.H{
			"error": "could not find user with such email",
		})
		return
	}
	user.Role = "Admin"
	user.UpdateUser()
	c.JSON(200, gin.H{
		"status": "user is now an admin",
		"errors": false,
	})
}

func CreateBorderlessVirtualAccount(c *gin.Context) {
	var input serializers.UserAccountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	kyc, err := models.GetKYCByUserID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KYC not found for user"})
		return
	}

	if kyc.BorderlessIdentityId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KYC not submitted to Borderless"})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	borderless := borderless.NewBorderless()
	userAccountId, err := models.GenerateAccountId()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountName := fmt.Sprintf("greybox-%s", userAccountId)

	accountResponse, err := borderless.CreateBorderlessAccount(accountName, kyc.BorderlessIdentityId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user is allowed to create a virtual account in their country
	if utils.IsBlockedCountry(user.CountryCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Virtual accounts are not available in user's country"})
		return
	}

	accountId := accountResponse["id"].(string)

	// Here we hardcode US as the country for the virtual accounts
	// this is because the virtual accounts only opperate in USA
	//
	// Above we already make the check to be sure the user is allowed to operate
	// a virtual account from their country of registration
	virtualAccountResponse, err := borderless.CreateBorderlessVirtualAccount(
		accountId,
		input.Fiat,
		input.Asset,
		"US",
		kyc.BorderlessIdentityId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	virtualAccountId := virtualAccountResponse["id"].(string)

	userAccount := models.UserAccounts{
		UserId:           userId,
		AccountId:        accountId,
		VirtualAccountId: virtualAccountId,
		Fiat:             input.Fiat,
		Asset:            input.Asset,
		Country:          user.Country,
	}

	if err := userAccount.CreateUserAccount(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Virtual account created successfully",
		"data":   userAccount,
	})
}

func GetUserAccounts(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accounts, err := models.GetUserAccountsByUserId(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Fetched user accounts successfully",
		"data":   accounts,
	})
}

func FilterUserAccounts(c *gin.Context) {
	var filters serializers.UserAccountsFilter
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts, err := models.FilterUserAccounts(filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Filtered user accounts successfully",
		"data":   accounts,
	})
}
