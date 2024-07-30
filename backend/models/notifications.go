package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hmac struct {
	*gorm.Model
	Secret uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"hmac_secret"`
}

func NewHmac() *Hmac {
	return &Hmac{}
}
