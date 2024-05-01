package models

import (
	"backend/utils/tokens"
	"errors"
	"fmt"
	"html"
	"math/rand"
	"net/mail"
	"strings"
	"unicode"

	"sync/atomic"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `gorm:"unique;not null" json:"email"`
	Password       string `json:"-"`
	Currency       string `json:"currency"`
	Country        string `json:"country"`
	Mnemonic       string `json:"mnemonic"`
	Xpub           string `json:"xpub"`
	AccountID      string `gorm:"unique" json:"account_id"`
	CustomerId     string `json:"customer_id"`
	AccountNumber  string `json:"account_number"`
	AccountCode    string `json:"account_code"`
	CountryCode    string `json:"country_code"`
	IsVerified     bool   `gorm:"default:false" json:"is_verified"`
	AccountAddress string `json:"account_address"`
	PrivateKey     string `json:"-"`
	CryptoCurrency string `gorm:"default:CELO" json:"crypto_currency"`

	UserImage string `json:"user_image"`
}

var counter uint64

// generateAccountCode generates a new account code
func GenerateAccountCode(prefix string) string {
	// Increment counter atomically and format it
	num := atomic.AddUint64(&counter, 1)
	return fmt.Sprintf("%02d_%s_%02d", num, prefix, num)
}

func GenerateAccountId() string {
	// Generate a random string of 6 characters
	random := randString(10)

	accountId := random

	// Check if account ID already exists
	if AlreadyExists(accountId) {
		// Generate a new random string and check again if needed
		return GenerateAccountId()
	}

	return accountId
}

func GenerateAccountNumber() string {
	// Generate a random string of 6 characters
	random := RandNumber(10)

	accountId := random

	// Check if account ID already exists
	if AccountNuberExists(accountId) {
		// Generate a new random string and check again if needed
		return GenerateAccountNumber()
	}

	return accountId
}

func randString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz1234567890"

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandNumber(n int) string {
	letters := "1234567890"

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func AlreadyExists(id string) bool {
	fmt.Println("id: ", id)
	var user User

	result := db.Where("account_id = ?", id).First(&user)

	return result.Error == nil

}

func AccountNuberExists(id string) bool {
	fmt.Println("id: ", id)
	var user User

	result := db.Where("account_number = ?", id).First(&user)

	return result.Error == nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func (u *User) SaveUser() error {
	//var err error
	err := db.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateUser() {
	db.Save(&u)
}

func (u *User) BeforeSave() error {
	if !ValidatePassword(u.Password) {
		return errors.New("password must have at least 8 characters, have at least a digit and at least an Upper case letter")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = CheckEmail(u.Email)

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

func CheckEmail(address string) string {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "invalid email adress"
	}
	return addr.Address
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
		return "", err
	}
	token, err := tokens.GenerateToken(u.ID)
	//fmt.Print("token: ", token)
	if err != nil {
		//fmt.Print("err: ", err)
		return "", err
	}

	return token, nil
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

	if err := db.Find(&user).Where("email = ?", lowercaseEmail).Error; err != nil {
		fmt.Println("user:", user)
		return user, false
	}

	return user, true
}
