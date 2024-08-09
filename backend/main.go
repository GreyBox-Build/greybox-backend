package main

import (
	"backend/controllers"
	"backend/middlewares"

	"backend/models"
	"log"
	"time"

	//"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORS() gin.HandlerFunc {
	// TO allow CORS
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
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
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//gin.SetMode(gin.ReleaseMode)

	db := models.InitializeDB()
	models.Migrate(db)
	r := gin.Default()

	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(CORS())

	r.Use(middlewares.AllowedHosts([]string{"localhost:8080", "34.227.150.136", "apis.greyboxpay.com", "wallet.greyboxpay.com"}))

	chains := r.Group("/api/v1/chains")
	{
		chains.GET("", controllers.FetchChain)
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
		//notification.Use(middlewares.JwtAuthMiddleware())
		notification.POST("/register-hmac", controllers.RegisterHmac)
	}

	master := r.Group("/api/v1")
	{
		//master.POST("/master-wallet", controllers.GenerateMasterWallet)
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
	}

	r.Run(":8080")
}
