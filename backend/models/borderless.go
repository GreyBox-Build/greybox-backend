package models

import "gorm.io/gorm"

type BorderlessRequest struct {
	*gorm.Model
	FiatAmount string `json:"fiatAmount"`
	Asset      string `json:"asset"`
	Country    string `json:"country"`
	UserId     uint   `json:"userId"`
	User       User   `gorm:"foreignKey:UserId" json:"user"`
	Status     string `json:"status"`
	TxId       string `json:"txId"`
	AccountId  string `json:"accountId"`
	TxHash     string `json:"txHash"`
	FeeAmount  string `json:"feeAmount"`
}


// CreateBorderlessRequest creates a new borderless request
func CreateBorderlessRequest(borderlessRequest *BorderlessRequest) error {
	return db.Create(borderlessRequest).Error
}