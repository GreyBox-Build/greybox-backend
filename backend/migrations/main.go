// cmd/migrate/main.go
package main

import (
	"backend/models"
	"backend/state"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if not in production
	state.LoadEnv()
	if state.AppConfig.AppEnv != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Load database configuration
	dbConfig := models.LoadDBConfigFromEnv()

	// 2. Initialize database connection
	db, err := models.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Perform database migrations
	err = models.Migrate(db,
		&models.User{},
		&models.Token{},
		&models.Transaction{},
		&models.XlmPublic{},
		&models.MasterWallet{},
		&models.WalletAddress{},
		&models.DepositRequest{},
		&models.WithdrawalRequest{},
		&models.HurupayRequest{},
		&models.BorderlessRequest{},
		&models.KYC{},
		&models.KYCData{},
		&models.UserAccounts{},
		&models.Bank{},
	)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
}
