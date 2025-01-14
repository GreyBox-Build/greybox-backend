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
	RequestId          string  `json:"request_id"`
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
	UserId          int32       `json:"user_id"`
	User            User        `gorm:"foreignKey:UserId" json:"user"`
	RequestId       string      `json:"request_id"`
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

func (t *Transaction) UpdateTransaction() error {
	return db.Save(t).Error
}
func GetTransactionByHash(hash, chain string) (*Transaction, error) {
	var transaction Transaction
	err := db.Preload("User").Where("hash = ? AND chain = ?", hash, chain).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func GetTransactionsByUserID(userId uint) ([]*Transaction, error) {
	var transactions []*Transaction

	err := db.Preload("User").Where("user_id = ?", userId).Find(&transactions).Error
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

func GetTransactionByRequestId(requestId string) (*Transaction, error) {
	var transaction Transaction
	err := db.Preload("User").Where("request_id = ?", requestId).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func GetHurupayRequest() ([]*HurupayRequest, error) {
	var hurupayRequests []*HurupayRequest
	err := db.Preload("User").Find(&hurupayRequests).Error
	if err != nil {
		return nil, err
	}
	return hurupayRequests, nil
}

func GetHurupayRequestById(id int) (*HurupayRequest, error) {
	var hurupayRequest HurupayRequest
	err := db.Preload("User").First(&hurupayRequest, id).Error
	if err != nil {
		return nil, err
	}
	return &hurupayRequest, nil
}

func GetHurupayStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	var totalRequests int64
	if err := db.Model(&HurupayRequest{}).Count(&totalRequests).Error; err != nil {
		return nil, err
	}
	stats["total_requests"] = totalRequests

	// Total successful requests with request type "on-ramp"
	var successfulOnRamp int64
	if err := db.Model(&HurupayRequest{}).
		Where("status = ? AND request_type = ?", "Completed", OnRamp).
		Count(&successfulOnRamp).Error; err != nil {
		return nil, err
	}
	stats["successful_on_ramp"] = successfulOnRamp

	// Total failed requests with request type "on-ramp"
	var failedOnRamp int64
	if err := db.Model(&HurupayRequest{}).
		Where("status = ? AND request_type = ?", "CANCELLED", OnRamp).
		Count(&failedOnRamp).Error; err != nil {
		return nil, err
	}
	stats["failed_on_ramp"] = failedOnRamp

	// Total pending requests for "off-ramp"
	var pendingOffRamp int64
	if err := db.Model(&HurupayRequest{}).
		Where("status = ? AND request_type = ?", "Pending", OffRamp).
		Count(&pendingOffRamp).Error; err != nil {
		return nil, err
	}
	stats["pending_off_ramp"] = pendingOffRamp

	var pendingOnRamp int64
	if err := db.Model(&HurupayRequest{}).
		Where("status = ? AND request_type = ?", "Pending", OnRamp).
		Count(&pendingOnRamp).Error; err != nil {
		return nil, err
	}
	stats["pending_on_ramp"] = pendingOnRamp
	// Total successful requests with request type "off-ramp"
	var successfulOffRamp int64
	if err := db.Model(&HurupayRequest{}).
		Where("status = ? AND request_type = ?", "Completed", OffRamp).
		Count(&successfulOffRamp).Error; err != nil {
		return nil, err
	}
	stats["successful_off_ramp"] = successfulOffRamp

	var failedOffRamp int64
	if err := db.Model(&HurupayRequest{}).
		Where("status = ? AND request_type = ?", "declined", OffRamp).
		Count(&failedOffRamp).Error; err != nil {
		return nil, err
	}
	stats["failed_off_ramp"] = failedOffRamp

	var CreatedOnramp int64

	if err := db.Model(&HurupayRequest{}).
		Where("status = ? AND request_type = ?", "CREATED", OnRamp).
		Count(&CreatedOnramp).Error; err != nil {
		return nil, err
	}
	stats["created_onramp"] = CreatedOnramp

	return stats, nil
}

type HurupayRequestSerializer struct {
	Amount             string      `json:"amount"`
	CountryCurrency    string      `json:"country_currency"`
	AccountNumber      string      `json:"account_number"`
	UserId             int32       `json:"user_id"`
	RequestId          string      `json:"request_id"`
	Status             string      `json:"status"`
	MobileNetwork      string      `json:"mobile_network"`
	ConfirmedAt        time.Time   `json:"confirmed_at"`
	CryptoChain        string      `json:"crypto_chain"`
	Token              string      `json:"token"`
	CountryCode        string      `json:"country_code"`
	MobileNumber       string      `json:"mobile_number"`
	RequestType        RequestType `json:"request_type"`
	UserFirstname      string      `json:"user_firstname"`
	UserLastname       string      `json:"user_lastname"`
	UserEmail          string      `json:"user_email"`
	UserAccountAddress string      `json:"user_account_address"`
	UserCryptoCurrency string      `json:"user_crypto_currency"`
	UserCountry        string      `json:"user_country_code"`
	Id                 uint        `json:"id"`
}

// ConvertToSerializer Function to convert HurupayRequest to HurupayRequestSerializer
func ConvertToSerializer(input interface{}) interface{} {
	switch v := input.(type) {
	case *HurupayRequest:
		return serializeRequest(*v)
	case []*HurupayRequest:
		var result []HurupayRequestSerializer
		for _, req := range v {
			result = append(result, serializeRequest(*req))
		}
		return result
	default:
		panic("Unsupported type")
	}
}

// Helper function to serialize a single HurupayRequest
func serializeRequest(request HurupayRequest) HurupayRequestSerializer {
	return HurupayRequestSerializer{
		Amount:             request.Amount,
		CountryCurrency:    request.CountryCurrency,
		AccountNumber:      request.AccountNumber,
		UserId:             request.UserId,
		RequestId:          request.RequestId,
		Status:             request.Status,
		MobileNetwork:      request.MobileNetwork,
		ConfirmedAt:        request.ConfirmedAt,
		CryptoChain:        request.CryptoChain,
		Token:              request.Token,
		CountryCode:        request.CountryCode,
		MobileNumber:       request.MobileNumber,
		RequestType:        request.RequestType,
		UserFirstname:      request.User.FirstName,
		UserLastname:       request.User.LastName,
		UserEmail:          request.User.Email,
		UserAccountAddress: request.User.AccountAddress,
		UserCryptoCurrency: request.User.CryptoCurrency,
		UserCountry:        request.User.CountryCode,
		Id:                 request.ID,
	}
}
