package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitializeDB() *gorm.DB {

	if os.Getenv("GIN_MODE") == "release" {

		host := os.Getenv("DB_HOST")

		port := os.Getenv("DB_PORT")

		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, user, password, dbname, port)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
		}

		log.Println("Connected to PostgreSQL database successfully!")
	} else {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect to SQLite database")
		}
		log.Println("connected to SQLite database successfully!")
	}

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Token{})
	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&XlmPublic{})
	db.AutoMigrate(&MasterWallet{})
	db.AutoMigrate(&WalletAddress{})
	db.AutoMigrate(&DepositRequest{})
	db.AutoMigrate(&WithdrawalRequest{})

}
