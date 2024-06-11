package models

import (
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID             uint   `gorm:"index" json:"user_id"`
	User               User   `gorm:"foreignKey:UserID" json:"user"`
	Amount             string `json:"amount"`
	Status             string `json:"status"`
	Chain              string `gorm:"default:celo" json:"chain"`
	Hash               string `json:"hash"`
	TransactionSubType string `json:"transaction_sub_type"`
	TransactionType    string `json:"transaction_type"`
	TransactionIndex   uint   `json:"transaction_index"`
	Address            string `json:"address"`
	BlockNumber        uint   `json:"block_number"`
	TransactionId      string `json:"transaction_id"`
	TransFee           string `json:"trans_fee"`
}

// WeiToGwei converts Wei to Gwei.
func WeiToGwei(wei *big.Int) *big.Int {
	gwei := new(big.Int).Div(wei, big.NewInt(1e9))
	return gwei
}

func CreateTransaction(userId uint) (*Transaction, error) {

	ref := uuid.New().String()
	result := db.Where("transaction_id = ?", ref).First(&Transaction{})

	if result.Error == nil {
		return CreateTransaction(userId)
	}

	transaction := &Transaction{
		UserID:           userId,
		Amount:           "0",
		Status:           "pending",
		Chain:            "celo",
		TransactionIndex: 1,
		BlockNumber:      0,
		TransactionId:    ref,
	}
	err := db.Create(transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *Transaction) SaveTransaction() error {
	return db.Create(t).Error
}
