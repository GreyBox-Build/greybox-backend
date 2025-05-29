package models

import (
	"backend/serializers"
	"backend/utils/tokens"
	"errors"
	"fmt"
	"html"
	"math/rand"
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
	if AccountNumberExists(accountId) {
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

	var user User

	result := db.Where("account_id = ?", id).First(&user)

	return result.Error == nil

}

func AccountNumberExists(id string) bool {
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

func (u *User) UpdateUserWithErrors() error {
	err := db.Model(&User{}).Where("id = ?", u.ID).Updates(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeSaveDetail() error {
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
		return "invalid email address"
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
	err := db.Create(&ua).Error
	if err != nil {
		return err
	}
	return nil
}

func (ua *UserAccounts) UpdateUserAccount() error {
	err := db.Save(&ua).Error
	if err != nil {
		return err
	}
	return nil
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
