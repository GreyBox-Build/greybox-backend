package models

import (
	"backend/state"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// Config holds database connection parameters.
// This makes configuration more explicit and testable.
type DBConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	TimeZone   string // Make TimeZone configurable or use UTC
	UseSQLite  bool   // Flag for development mode
	SQLitePath string
}

// NewDB initializes and returns a GORM database connection.
// It returns (*gorm.DB, error) to allow for proper error handling by the caller.
func NewDB(cfg DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if !cfg.UseSQLite { // Changed from GIN_MODE to a config flag
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.TimeZone)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
		}
		log.Println("Connected to PostgreSQL database successfully!")
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
		}
		log.Println("Connected to SQLite database successfully!")
	}

	return db, nil
}

// LoadDBConfigFromEnv loads database configuration from environment variables.
// This function would typically be in a 'config' package or main.
func LoadDBConfigFromEnv() DBConfig {
	isRelease := state.AppConfig.GinMode == "release"

	cfg := DBConfig{
		UseSQLite:  !isRelease, // If not release, use SQLite
		SQLitePath: "test.db",  // Default for SQLite
		TimeZone:   "UTC",      // Default to UTC
	}

	if isRelease {
		cfg.Host = state.AppConfig.DBHost
		cfg.Port = state.AppConfig.DBPort
		cfg.User = state.AppConfig.DBUser
		cfg.Password = state.AppConfig.DBPassword
		cfg.DBName = state.AppConfig.DBName
		// Optional: if you still want a specific timezone from env
		if tz := state.AppConfig.DBTimezone; tz != "" {
			cfg.TimeZone = tz
		}
	}
	// Add validation for required fields if not using SQLite
	if !cfg.UseSQLite && (cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Password == "" || cfg.DBName == "") {
		log.Fatal("Missing one or more required database environment variables for release mode.")
	}
	return cfg
}

// Migrate performs database migrations for the given models.
// It takes a list of models to migrate, making it more flexible.
func Migrate(db *gorm.DB, modelsToMigrate ...interface{}) error {
	log.Println("Starting database migrations...")
	err := db.AutoMigrate(modelsToMigrate...)
	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	log.Println("Database migrations completed successfully!")
	return nil
}

// InitializeDB sets up the database connection and performs migrations.
// It returns the *gorm.DB instance, which is crucial for the rest of your application.
func InitializeDB() *gorm.DB {
	// 1. Load configuration
	dbConfig := LoadDBConfigFromEnv()

	// 2. Initialize database connection
	db, err = NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	log.Println("Database successfully connected!")
	return db
}
