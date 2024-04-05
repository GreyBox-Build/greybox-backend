package main

import (
	"backend/models"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	db := models.InitializeDB()
	models.Migrate(db)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))
}
