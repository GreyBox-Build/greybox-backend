package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID             uint    `gorm:"index" json:"user_id"`
	User               User    `gorm:"foreignKey:UserID" json:"user"`
	Amount             string  `json:"amount"`
	Status             string  `json:"status"`
	Chain              string  `gorm:"default:celo" json:"chain"`
	Hash               string  `json:"hash"`
	TransactionSubType string  `json:"transaction_sub_type"`
	TransactionType    string  `json:"transaction_type"`
	TransactionIndex   uint    `json:"transaction_index"`
	Address            string  `json:"address"`
	BlockNumber        uint    `json:"block_number"`
	TransactionId      string  `json:"transaction_id"`
	TransFee           float64 `json:"trans_fee"`
	Description        string  `json:"description"`
	CounterAddress     string  `json:"counter_address"`
	TokenId            *string `json:"token_id"`
	Asset              string  `json:"asset"`
}

type DepositRequest struct {
	gorm.Model
	UserID          uint      `gorm:"index" json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"user"`
	Status          string    `json:"status"`
	Ref             string    `json:"ref"`
	CountryCode     string    `json:"country_code"`
	DepositBank     string    `json:"deposit_bank"`
	AccountNumber   string    `json:"account_number"`
	ConfirmedAt     time.Time `json:"confirmed_at"`
	VerifiedById    uint      `json:"verified_by_id"`
	VerifiedBy      User      `gorm:"foreignKey:VerifiedById" json:"verified_by"`
	Currency        string    `json:"currency"`
	FiatAmount      string    `json:"fiat_amount"`
	ProposedAsset   string    `json:"proposed_asset"`
	AccountName     string    `json:"account_name"`
	AssetEquivalent string    `json:"asset_equivalent"`
}

type WithdrawalRequest struct {
	gorm.Model
	UserID        uint      `gorm:"index" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	Status        string    `json:"status"`
	CryptoAmount  string    `json:"amount"`
	Chain         string    `json:"chain"`
	Hash          string    `json:"hash"`
	Address       string    `json:"address"`
	BankName      string    `json:"bank_name"`
	AccountName   string    `json:"account_name"`
	AccountNumber string    `json:"account_number"`
	ConfirmedAt   time.Time `json:"confirmed_at"`
	VerifiedById  uint      `json:"verified_by_id"`
	VerifiedBy    User      `gorm:"foreignKey:VerifiedById" json:"verified_by"`
}

func (d *DepositRequest) SaveDepositRequest() error {
	return db.Create(d).Error
}

func (d *DepositRequest) UpdateDepositRequest() error {
	return db.Save(d).Error
}

func (w *WithdrawalRequest) SaveWithdrawalRequest() error {
	return db.Create(w).Error
}

func (w *WithdrawalRequest) UpdateWithdrawalRequest() error {
	return db.Save(w).Error
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

func GetTransactionByHash(hash, chain string) (*Transaction, error) {
	var transaction Transaction
	err := db.Where("hash = ? AND chain = ?", hash, chain).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func GetTransactionsByUserID(userId uint, chain string) ([]*Transaction, error) {
	var transactions []*Transaction
	err := db.Where("user_id = ? AND chain = ?", userId, chain).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func GenerateRequestReference() string {
	return uuid.New().String()
}
