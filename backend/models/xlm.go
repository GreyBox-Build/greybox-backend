package models

import (
	"errors"

	"gorm.io/gorm"
)

type XlmPublic struct {
	gorm.Model
	Xpub   string `json:"xpub"`
	Secret string `json:"secret"`
}

func GetLatestXlmPublic() (*XlmPublic, error) {
	var x XlmPublic
	err := db.Order("id DESC").First(&x).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &x, errors.New("no xpub found")
		}
		return &x, err
	}
	return &x, nil
}

func (w *XlmPublic) Save() error {
	return db.Create(w).Error
}
