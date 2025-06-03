package models

import (
	"backend/serializers"
	"backend/state"
	"backend/utils/tokens"
	"crypto/rand"
	"errors"
	"fmt"
	"html"
	"log"
	"math/big"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"sync/atomic"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Email           string         `gorm:"unique;not null" json:"email"`
	Password        string         `json:"-"`
	Currency        string         `json:"currency"`
	Country         string         `json:"country"`
	Mnemonic        string         `json:"-"`
	Xpub            string         `json:"-"`
	CountryCode     string         `json:"country_code"`
	IsVerified      bool           `gorm:"default:false" json:"is_verified"`
	AccountAddress  string         `json:"account_address"`
	PrivateKey      string         `json:"-"`
	CryptoCurrency  string         `gorm:"default:CELO" json:"crypto_currency"`
	UserImage       string         `json:"user_image"`
	SignatureId     string         `json:"-"`
	TokenAddress    string         `json:"token_address"`
	Index           uint64         `json:"-"`
	PreviousBalance float32        `json:"-"`
	Role            string         `gorm:"default:Customer" json:"role"`
	UserAccounts    []UserAccounts `gorm:"foreignKey:UserId" json:"user_accounts"`
}

type UserAccounts struct {
	gorm.Model
	UserId           uint      `json:"user_id"`
	AccountId        string    `json:"account_id"`
	Asset            string    `json:"asset"`
	Fiat             string    `json:"fiat"`
	Country          string    `json:"country"`
	VirtualAccountId string    `json:"virtual_account_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

const maxRetry = 5

// generateAccountCode generates a new account code
func GenerateAccountCode(prefix string) string {
	// Increment counter atomically and format it
	num := atomic.AddUint64(&state.AccountCodeCounter, 1)
	return fmt.Sprintf("%02d_%s_%02d", num, prefix, num)
}

func GenerateAccountId() (string, error) {
	for range maxRetry {
		id, err := randString(10)
		if err != nil {
			return "", err
		}
		if !AlreadyExists(id) {
			return id, nil
		}
	}
	return "", errors.New("failed to generate unique account ID after multiple attempts")
}

func GenerateAccountNumber() (string, error) {
	for range maxRetry {
		num, err := randNumber(10)
		if err != nil {
			return "", err
		}
		if !AccountNumberExists(num) {
			return num, nil
		}
	}
	return "", errors.New("failed to generate unique account number after multiple attempts")
}

func secureRandomString(charset string, n int) (string, error) {
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("secure random generation failed: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func randString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz1234567890"
	return secureRandomString(letters, n)
}

func randNumber(n int) (string, error) {
	const numbers = "1234567890"
	return secureRandomString(numbers, n)
}

func CheckEmail(address string) (string, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", errors.New("invalid email address format")
	}
	return addr.Address, nil
}

func ValidatePassword(password string) bool {
	// Check length
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter
	hasUppercase := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
			break
		}
	}

	if !hasUppercase {
		return false
	}

	// Check for at least one digit
	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}

	return hasDigit
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(email, password string) (string, error) {
	u := User{}
	//var err error

	err := db.Model(User{}).Where("email=?", email).Take(&u).Error
	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", errors.New("Invalid password")
	}
	token, err := tokens.GenerateToken(u.ID)
	//fmt.Print("token: ", token)
	if err != nil {
		//fmt.Print("err: ", err)
		return "", err
	}

	return token, nil
}

func AlreadyExists(id string) bool {

	var user User

	result := db.Where("account_id = ?", id).First(&user)

	return result.Error == nil

}

func AccountNumberExists(id string) bool {
	log.Println("id: ", id)
	var user User

	result := db.Where("account_number = ?", id).First(&user)

	return result.Error == nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func (u *User) SaveUser() error {
	if err := u.SanitizeAndValidate(); err != nil {
		return err
	}
	if err := u.HashPassword(); err != nil {
		return err
	}
	return db.Create(&u).Error
}

func (u *User) UpdateUser() {
	db.Save(&u)
}

func (u *User) UpdateUserWithErrors() error {
	err := db.Model(&User{}).Where("id = ?", u.ID).Updates(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SanitizeAndValidate() error {
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))

	email, err := CheckEmail(u.Email)
	if err != nil {
		return err
	}
	u.Email = email

	if !ValidatePassword(u.Password) {
		return errors.New("password must be at least 8 characters long, contain a digit, and an uppercase letter")
	}

	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := db.First(&u, uid).Error; err != nil {
		return u, errors.New("user not found")
	}

	u.PrepareGive()

	return u, nil

}

func DeleteUserByID(userID uint) error {
	// Find the user by ID
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}

	// Delete the user
	result = db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FindUserByEmail(email string) (User, bool) {
	var user User

	lowercaseEmail := strings.ToLower(email)

	if err := db.Where("email = ?", lowercaseEmail).First(&user).Error; err != nil {

		return user, false
	}

	return user, true
}

func FindUserByAddress(address string) (User, error) {
	var user User
	if err := db.Where("account_address = ?", address).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func FindAdmins() ([]User, error) {
	var users []User
	if err := db.Where("role = ?", "Admin").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ua *UserAccounts) BeforeCreate(tx *gorm.DB) (err error) {
	ua.CreatedAt = time.Now()
	return
}

func (ua *UserAccounts) BeforeSave(tx *gorm.DB) (err error) {
	ua.UpdatedAt = time.Now()
	return
}

func (ua *UserAccounts) CreateUserAccount() error {
	return db.Create(&ua).Error
}

func (ua *UserAccounts) UpdateUserAccount() error {
	return db.Save(&ua).Error
}

func GetUserAccountsByUserId(userId uint) ([]UserAccounts, error) {
	var userAccounts []UserAccounts
	err := db.Where("user_id = ?", userId).Find(&userAccounts).Error
	if err != nil {
		return nil, err
	}
	return userAccounts, nil
}

func GetUserAccountById(id uint) (UserAccounts, error) {
	var userAccount UserAccounts
	err := db.First(&userAccount, id).Error
	if err != nil {
		return userAccount, err
	}
	return userAccount, nil
}

func FilterUserAccounts(filter serializers.UserAccountsFilter) ([]UserAccounts, error) {
	var userAccounts []UserAccounts
	query := db

	// Build the query dynamically based on provided filters
	if filter.UserId != nil {
		query = query.Where("user_id = ?", filter.UserId)
	}

	if filter.Asset != nil {
		query = query.Where("asset = ?", filter.Asset)
	}

	if filter.Country != nil {
		query = query.Where("country = ?", filter.Country)
	}

	if filter.Fiat != nil {
		query = query.Where("fiat = ?", filter.Fiat)
	}

	// Execute the query
	err := query.Find(&userAccounts).Error
	return userAccounts, err
}
