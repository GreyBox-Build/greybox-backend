package models

import (
	"backend/utils/tokens"
	"errors"
	"html"
	"net/mail"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `json:"-"`
	Currency  string `json:"currency"`
	Country   string `json:"country"`

	UserImage string `json:"user_image"`
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func (u *User) SaveUser() (*User, error) {
	//var err error
	err := db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
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
