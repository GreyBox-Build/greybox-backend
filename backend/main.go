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
	gin.SetMode(gin.ReleaseMode)

	db := models.InitializeDB()
	models.Migrate(db)
	r := gin.Default()

	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(CORS())

	chains := r.Group("/api/v1/chains")
	{
		chains.GET("", controllers.FetchChain)
	}

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
	}

	trans := r.Group("/api/v1/transaction")
	{
		trans.Use(middlewares.JwtAuthMiddleware())
		trans.GET("/on-ramp", controllers.RetrieveOnRampParamsV1)
		trans.GET("", controllers.GetUserTransactions)
		trans.GET("/hash", controllers.GetTransactionsByHash)
		trans.POST("/off-ramp", controllers.OffRampTransaction)
		trans.POST("/sign-url", controllers.SignUrl)
		trans.GET("/verify/:transaction_id", controllers.KMStransactionVerification)
	}

	notification := r.Group("/api/v1/notification")
	{
		//notification.Use(middlewares.JwtAuthMiddleware())
		notification.POST("/register-hmac", controllers.RegisterHmac)
	}

	r.Run(":8080")
}
