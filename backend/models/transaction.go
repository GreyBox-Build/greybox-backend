package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestType string

// Define constants for the possible values of the enum
const (
	OnRamp  RequestType = "On-ramp"
	OffRamp RequestType = "Off-ramp"
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
	UserID          uint      `json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"user"`
	Status          string    `json:"status"`
	Ref             string    `json:"ref"`
	CountryCode     string    `json:"country_code"`
	DepositBank     string    `json:"deposit_bank"`
	AccountNumber   string    `json:"account_number"`
	ConfirmedAt     time.Time `json:"confirmed_at"`
	VerifiedById    *uint     `gorm:"default:null" json:"verified_by_id"`
	VerifiedBy      *User     `gorm:"foreignKey:VerifiedById" json:"verified_by"`
	Currency        string    `json:"currency"`
	FiatAmount      string    `json:"fiat_amount"`
	ProposedAsset   string    `json:"proposed_asset"`
	AccountName     string    `json:"account_name"`
	AssetEquivalent string    `json:"asset_equivalent"`
}

type WithdrawalRequest struct {
	gorm.Model
	UserID         uint      `gorm:"index" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user"`
	Status         string    `json:"status"`
	CryptoAmount   string    `json:"crypto_amount"`
	Chain          string    `json:"chain"`
	Hash           string    `json:"hash"`
	Address        string    `json:"address"`
	BankName       string    `json:"bank_name"`
	AccountName    string    `json:"account_name"`
	AccountNumber  string    `json:"account_number"`
	ConfirmedAt    time.Time `json:"confirmed_at"`
	VerifiedById   *uint     `gorm:"default:null" json:"verified_by_id"`
	VerifiedBy     *User     `gorm:"foreignKey:VerifiedById" json:"verified_by"`
	BankRef        string    `json:"bank_ref"`
	Asset          string    `json:"asset"`
	EquivalentFiat string    `json:"equivalent_fiat"`
	FiatCurrency   string    `json:"fiat_currency"`
}

type HurupayRequest struct {
	gorm.Model
	Amount          string      `json:"amount"`
	CountryCurrency string      `json:"country_currency"`
	AccountNumber   string      `json:"account_number"`
	UserId          int32       `gorm:"index" json:"user_id"`
	User            User        `gorm:"foreignKey:UserId" json:"user"`
	RequestId       string      `gorm:"index" json:"request_id"`
	Status          string      `json:"status"`
	MobileNetwork   string      `json:"mobile_network"`
	ConfirmedAt     time.Time   `json:"confirmed_at"`
	CryptoChain     string      `json:"crypto_chain"`
	Token           string      `json:"token"`
	CountryCode     string      `json:"country_code"`
	MobileNumber    string      `json:"mobile_number"`
	RequestType     RequestType `json:"request_type"`
}

func (h *HurupayRequest) SaveHurupayRequest() error {
	return db.Create(h).Error
}

func (h *HurupayRequest) UpdateHurupayRequest() error {
	return db.Save(h).Error
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
	err := db.Preload("User").Where("hash = ? AND chain = ?", hash, chain).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func GetTransactionsByUserID(userId uint, chain string) ([]*Transaction, error) {
	var transactions []*Transaction
	err := db.Preload("User").Where("user_id = ? AND chain = ?", userId, chain).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func GenerateRequestReference() string {
	return uuid.New().String()
}

func FilterDepositRequests(ref, currency, fiatAmount, accountNumber, status, countryCode, cryptoAsset string) ([]DepositRequest, error) {
	var depositRequests []DepositRequest

	query := db.Model(&DepositRequest{})

	if ref != "" {
		query = query.Where("ref = ?", ref)
	}
	if currency != "" {
		query = query.Where("currency = ?", currency)
	}
	if fiatAmount != "" {
		query = query.Where("fiat_amount = ?", fiatAmount)
	}
	if accountNumber != "" {
		query = query.Where("account_number = ?", accountNumber)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if countryCode != "" {
		query = query.Where("country_code = ?", countryCode)
	}
	if cryptoAsset != "" {
		query = query.Where("proposed_asset = ?", cryptoAsset)
	}

	err := query.Preload("User").Find(&depositRequests).Error
	if err != nil {
		return nil, err
	}

	return depositRequests, nil
}

func FilterWithdrawalRequests(status, chain, hash, address, accountNumber string) ([]WithdrawalRequest, error) {
	var withdrawalRequests []WithdrawalRequest

	query := db.Model(&WithdrawalRequest{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if chain != "" {
		query = query.Where("chain = ?", chain)
	}
	if hash != "" {
		query = query.Where("hash = ?", hash)
	}
	if address != "" {
		query = query.Where("address = ?", address)
	}
	if accountNumber != "" {
		query = query.Where("account_number = ?", accountNumber)
	}

	err := query.Preload("User").Find(&withdrawalRequests).Error
	if err != nil {
		return nil, err
	}

	return withdrawalRequests, nil
}

func GetDepositRequest(id int) (*DepositRequest, error) {
	var depositRequest DepositRequest
	err := db.Preload("User").First(&depositRequest, id).Error
	if err != nil {
		return nil, err
	}
	return &depositRequest, nil
}

func GetWithdrawalRequest(id int) (*WithdrawalRequest, error) {
	var withdrawalRequest WithdrawalRequest
	err := db.Preload("User").First(&withdrawalRequest, id).Error
	if err != nil {
		return nil, err
	}
	return &withdrawalRequest, nil
}

func GetHurupayRequestRequestId(requestId string) (*HurupayRequest, error) {
	var hurupayRequest HurupayRequest
	err := db.Preload("User").First(&hurupayRequest, "request_id = ?", requestId).Error
	if err != nil {
		return nil, err
	}
	return &hurupayRequest, nil
}
