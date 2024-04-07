package models

import (
	"time"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	Token      string    `json:"token"`
	ExpireAt   time.Time `json:"expired_at"`
	HasExpired bool      `gorm:"default:false" json:"has_expired"`
}

func SaveToken(token *Token) error {
	return db.Create(&token).Error
}

func GenerateRecoveryToken(userID uint) (Token, error) {
	token := Token{
		UserID:   userID,
		Token:    uuid.NewString(),
		ExpireAt: time.Now().Add(time.Hour), // Token expires in 1 hour
	}

	if err := SaveToken(&token); err != nil {
		return token, err
	}
	return token, nil
}

func CheckTokenValid(token string) (Token, error) {
	var Rectoken Token

	if err := db.Where("token = ?", token).Where("has_expired = ?", false).First(&Rectoken).Error; err != nil {
		return Rectoken, err
	}

	if Rectoken.ExpireAt.Before(time.Now()) {
		return Rectoken, errors.New("the token is invalid i.e it has expired")
	}

	return Rectoken, nil
}

func (u *Token) UpdateToken() error {

	db.Save(&u)
	return nil
}
