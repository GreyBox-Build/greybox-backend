package models

import (
	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	Name          string  `gorm:"type:varchar(255);not null;unique" json:"name"`
	Country       string  `gorm:"type:varchar(100);not null" json:"country"`
	RoutingNumber *string `gorm:"type:varchar(20);null" json:"routing_number"`
	SwiftCode     *string `gorm:"type:varchar(20);null" json:"swift_code"`
	State         *string `gorm:"type:varchar(100)" json:"state"`
	Street        *string `gorm:"type:varchar(255)" json:"street"`
	ZipCode       *string `gorm:"type:varchar(20)" json:"zip_code"`
	City          *string `gorm:"type:varchar(100)" json:"city"`
}

func FilterBanks(country string) ([]Bank, error) {
	var banks []Bank

	result := db.Where("country = ?", country).Find(&banks)

	if result.Error != nil {
		return nil, result.Error
	}

	return banks, nil
}

func GetBankData(id int) (*Bank, error) {
	var bank Bank

	result := db.First(&bank, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &bank, nil
}
