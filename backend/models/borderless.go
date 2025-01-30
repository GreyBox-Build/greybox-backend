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

// GetBorderlessRequestByTxId fetches a borderless request by txId
func GetBorderlessRequestByTxId(txId string) (*BorderlessRequest, error) {
	var borderlessRequest BorderlessRequest
	err := db.Preload("User").Where("tx_id = ?", txId).First(&borderlessRequest).Error
	return &borderlessRequest, err
}

// UpdateBorderlessRequest updates a borderless request
func (borderlessRequest *BorderlessRequest) UpdateBorderlessRequest() error {
	return db.Save(borderlessRequest).Error
}

