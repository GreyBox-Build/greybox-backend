package main

import (
	"backend/controllers"
	"backend/middlewares"
	"backend/models"
	"backend/state"
	"time"

	//"github.com/gin-contrib/cors"

	"slices"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	// TO allow CORS
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "localhost:8080", "34.227.150.136", "apis.greyboxpay.com", "wallet.greyboxpay.com", "https://wallet.greyboxpay.com", "https://apis.greyboxpay.com"} // Add other allowed origins if needed

		if slices.Contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	lagos, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		// Handle error if time zone loading fails
		panic(err)
	}
	time.Local = lagos

	// load env
	state.LoadEnv()

	//gin.SetMode(gin.ReleaseMode)

	db := models.InitializeDB()
	models.Migrate(db)
	r := gin.Default()

	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(CORS())

	r.Use(middlewares.AllowedHosts([]string{"localhost:3000", "http://localhost:3000", "localhost:8080", "34.227.150.136", "apis.greyboxpay.com", "wallet.greyboxpay.com"}))

	chains := r.Group("/api/v1/chains")
	{
		chains.GET("", controllers.FetchChain)
	}

	networks := r.Group("/api/v1/networks")
	{
		networks.GET("", controllers.FetchNetwork)
	}

	XchangeRate := r.Group("/api/v1/exchange-rate")
	{
		XchangeRate.GET("", controllers.GetExchangeRate)
	}

	//r.Use(middlewares.JwtAuthMiddleware())

	public := r.Group("/api/v1/user")
	{
		public.POST("/login", controllers.FetchAuthenticatedUserToken)
		public.POST("/forget-password", controllers.ForgetPassword)
		public.POST("/reset-password", controllers.ResetPassword)
	}

	publicV2 := r.Group("/api/v2/user")
	{
		publicV2.POST("/register", controllers.CreateAccountV2)
		publicV2.Use(middlewares.JwtAuthMiddleware()).POST("/account", controllers.CreateBorderlessVirtualAccount)
		publicV2.Use(middlewares.JwtAuthMiddleware()).GET("/account", controllers.GetUserAccounts)
		publicV2.Use(middlewares.JwtAuthMiddleware()).Use(middlewares.IsAdmin()).GET("/accounts", controllers.FilterUserAccounts)
	}

	kyc := r.Group("/api/v2/kyc")
	{
		kyc.Use(middlewares.JwtAuthMiddleware())
		kyc.GET("/mine", controllers.GetUserKYC)
		kyc.POST("", controllers.CreateKYC)
		kyc.PATCH("", controllers.UpdateKYC)
		kyc.DELETE("/:id", controllers.DeleteKYC)
		kyc.Use(middlewares.IsAdmin()).GET("", controllers.GetKYCS)
		kyc.Use(middlewares.IsAdmin()).PATCH("/:id/approve", controllers.ApproveKYC)
		kyc.Use(middlewares.IsAdmin()).PATCH("/:id/reject", controllers.RejectKYC)
	}

	user := r.Group("/api/v1/auth")
	{
		user.Use(middlewares.JwtAuthMiddleware())
		user.GET("/user", controllers.GetAuthenticatedUser)
		user.POST("/make-admin", controllers.MakeAdmin)
	}

	trans := r.Group("/api/v1/transaction")
	{
		trans.Use(middlewares.JwtAuthMiddleware())
		trans.GET("/on-ramp", controllers.RetrieveOnRampParamsV1)
		trans.GET("", controllers.GetUserTransactions)
		trans.GET("/hash", controllers.GetTransactionsByHash)
		trans.POST("/off-ramp", controllers.OffRampTransaction)
		trans.POST("/sign-url", controllers.SignUrl)
	}

	notification := r.Group("/api/v1/notification")
	{
		//notification.Use(middlewares.WebhookSignatureMiddleware())
		notification.POST("/on-ramp", controllers.OnRampNotification)
		notification.POST("/off-ramp", controllers.OffRampNotification)

		//notification.POST("/register-hmac", controllers.RegisterHmac)
	}

	master := r.Group("/api/v1")
	{
		//master.POST("/master-wallet", controllers.GenerateMasterWallet)
		master.GET("/master-wallets", controllers.GetMasterWallets)
		master.GET("/master-wallet", controllers.GetMasterWallet)

	}

	transV2 := r.Group("/api/v2/transaction")
	{
		transV2.Use(middlewares.JwtAuthMiddleware())
		transV2.GET("/equivalent-amount", controllers.AmountToReceive)
		transV2.GET("/destination-bank", controllers.GetDestinationBankAccount)
		transV2.GET("/reference", controllers.GenerateReference)
		transV2.POST("/on-ramp", controllers.OnRampV2)
		transV2.POST("/off-ramp", controllers.OffRampV2)
		transV2.GET("/on-ramp/mobile/equivalent-amount", controllers.MobileMoneyAmountToReceive)
		transV2.POST("/on-ramp/mobile", controllers.MobileMoneyOnRamp)
		transV2.POST("/off-ramp/mobile", controllers.MobileMoneyOffRamp)

	}

	requests := r.Group("/api/v1/requests")
	{
		requests.Use(middlewares.JwtAuthMiddleware())
		requests.Use(middlewares.IsAdmin())
		requests.GET("/on-ramp", controllers.FetchOnRampRequests)
		requests.GET("/off-ramp", controllers.FetchOffRampRequests)
		requests.GET("/on-ramp/:id", controllers.GetOnRampRequest)
		requests.GET("/off-ramp/:id", controllers.GetOffRampRequest)
		requests.POST("/on-ramp/:id/verify", controllers.VerifyOnRamp)
		requests.POST("/off-ramp/:id/verify", controllers.VerifyOffRamp)
		requests.GET("/hurupay-requests", controllers.ListHurupayRequest)
		requests.GET("/hurupay-requests/:id", controllers.GetHurupayRequest)
		requests.GET("/hurupay-requests/stats", controllers.GetHurupayStats)
	}

	payments := r.Group("/api/v1/payments")
	{
		payments.Use(middlewares.JwtAuthMiddleware())
		payments.GET("/banks", controllers.FilterBank)
		payments.POST("/borderless-onramp", controllers.BorderLessOnramp)
		payments.POST("/borderless-offramp", controllers.BorderLessOffRamp)
		payments.POST("/borderless-onramp/mobilemoney", controllers.BorderlessMobileMoneyOnRamp)
		payments.POST("/borderless-offramp/mobilemoney", controllers.BorderlessMobileMoneyOffRamp)
	}

	webhook := r.Group("/api/v1/webhook")
	webhook.Use(middlewares.SignatureMiddleware("webhook_rsa"))
	{
		webhook.POST("/borderless", controllers.BorderlessNotification)
	}

	r.Run(":8080")
}
