package state

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

type Config struct {
	// App Config
	AppEnv   string
	AdminKey string
	GinMode  string

	// Database Config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBTimezone string

	// Mail Config
	EmailUser     string
	EmailPassword string

	// Borderless Config
	BorderlessClientId         string
	BorderlessClientSecret     string
	BorderlessAccountId        string
	BorderlessBaseUrl          string
	BorderlessBusinessIdentity string

	// Tatum Config
	TatumTestApiKey       string
	TatumWebhookUrl       string
	TatumSubscriptionType string
	TatumBaseUrl          string

	// Hurupay Config
	HurupayApiKey string

	// Moonpay Config
	MoonpayTestApiKey string

	// Celo Config
	SourceParam   string
	XClientId     string
	XClientSecret string

	// PrivateKey Encryption
	EncryptionKey string

	// Other Config
	HmacSecret string

	// Jwt Config
	ApiSecret                string
	TokenExpirationInMinutes int

	// PasswordReset
	PasswordResetLink string
}

var AppConfig *Config

// ApiSecret for Jwt signing and validation
var ApiSecret []byte

// Cache for temporarily holding borderless related data
var BorderlessCache = cache.New(60*time.Minute, 100*time.Minute)

// AccountCode Counter
var AccountCodeCounter uint64

// Function to Load Environment Variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// all configs using os.Getenv are optional
	// mustGetEnv makes sure required envs are loaded
	AppConfig = &Config{
		AppEnv:                     os.Getenv("APP_ENV"),
		AdminKey:                   mustGetEnv("ADMIN_KEY"),
		GinMode:                    mustGetEnv("GIN_MODE"),
		DBHost:                     mustGetEnv("DB_HOST"),
		DBPort:                     mustGetEnv("DB_PORT"),
		DBUser:                     mustGetEnv("DB_USER"),
		DBPassword:                 mustGetEnv("DB_PASSWORD"),
		DBName:                     mustGetEnv("DB_NAME"),
		DBTimezone:                 os.Getenv("DB_TIMEZONE"),
		EmailUser:                  mustGetEnv("EMAIL_USER"),
		EmailPassword:              mustGetEnv("EMAIL_PASSWORD"),
		BorderlessClientId:         mustGetEnv("BORDERLESS_CLIENT_ID"),
		BorderlessClientSecret:     mustGetEnv("BORDERLESS_CLIENT_SECRET"),
		BorderlessAccountId:        mustGetEnv("BORDERLESS_ACCOUNT_ID"),
		BorderlessBaseUrl:          mustGetEnv("BORDERLESS_BASE_URL"),
		BorderlessBusinessIdentity: mustGetEnv("BORDERLESS_BUSINESS_IDENTITY"),
		TatumTestApiKey:            mustGetEnv("TATUM_API_KEY_TEST"),
		TatumWebhookUrl:            os.Getenv("WEBHOOK_URL"),
		TatumSubscriptionType:      mustGetEnv("SUBSCRIPTION_TYPE"),
		TatumBaseUrl:               mustGetEnv("TATUM_BASE_URL"),
		HurupayApiKey:              mustGetEnv("HURUPAY_API_KEY"),
		MoonpayTestApiKey:          os.Getenv("MOONPAY_API_KEY_TEST"),
		SourceParam:                mustGetEnv("SOURCE_PARAM"),
		XClientId:                  mustGetEnv("X_CLIENT_ID"),
		XClientSecret:              mustGetEnv("X_CLIENT_SECRET"),
		HmacSecret:                 os.Getenv("HMAC_SECRET"),
		ApiSecret:                  mustGetEnv("API_SECRET"),
		TokenExpirationInMinutes:   mustGetEnvAsInt("TOKEN_EXPIRATION_IN_MINUTES"),
		EncryptionKey:              mustGetEnv("ENCRYPTION_KEY"),
		PasswordResetLink:          mustGetEnv("PASSWORD_RESET_LINK"),
	}

	ApiSecret = []byte(AppConfig.ApiSecret)
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}

func mustGetEnvAsInt(key string) int {
	value := os.Getenv(key)
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid %s value. Error during conversion: %v", key, err)
	}

	if parsedValue < 1 {
		log.Fatalf("Invalid %s value, cant be less than 1 minute", key)
	}

	return parsedValue
}
