package serializers

import "time"

type Chain struct {
	Celo    string `json:"celo"`
	Stellar string `json:"stellar"`
	Polygon string `json:"polygon"`
}

const (
	ChainCelo    = "CELO"
	ChainStellar = "XLM"
	ChainPolygon = "MATIC"
)

var Chains = Chain{
	Celo:    ChainCelo,
	Stellar: ChainStellar,
	Polygon: ChainPolygon,
}

type User struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
	Country     string `json:"country" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
	Chain       string `json:"chain"`
}

type LoginSerializer struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Data struct {
	Chain          string `json:"chain" binding:"required"`
	Name           string `json:"name" binding:"required"`
	ImageURL       string `json:"image_url" binding:"required"`
	AvailableToken string `json:"availableToken" binding:"required"`
}

type Account struct {
	ID            string `json:"id"`
	AccountID     string `json:"account_id"`
	Sequence      string `json:"sequence"`
	SubentryCount int    `json:"subentry_count"`
	Balances      []struct {
		Balance            string `json:"balance"`
		Limit              string `json:"limit"`
		BuyingLiabilities  string `json:"buying_liabilities"`
		SellingLiabilities string `json:"selling_liabilities"`
		AssetType          string `json:"asset_type"`
		AssetCode          string `json:"asset_code"`
		AssetIssuer        string `json:"asset_issuer"`
	} `json:"balances"`
	Signers []map[string]interface{} `json:"signers"`
	Data    map[string]interface{}   `json:"data"`
}

type AdminForm struct {
	Key       string `json:"key" binding:"required"`
	UserEmail string `json:"user_email" binding:"required,email"`
}

type NetworkData struct {
	CountryCode  string   `json:"countryCode" binding:"required"`
	CountryName  string   `json:"countryName" binding:"required"`
	CurrencyCode string   `json:"currencyCode" binding:"required"`
	Networks     []string `json:"networks" binding:"required"`
	MobileCode   string   `json:"mobileCode" binding:"required"`
}

type UserAccountRequest struct {
	Asset string `json:"asset" binding:"required"`
	Fiat  string `json:"fiat" binding:"required"`
}

type UserAccounts struct {
	AccountId        string    `json:"account_id"`
	Asset            string    `json:"asset"`
	Fiat             string    `json:"fiat"`
	Country          string    `json:"country"`
	VirtualAccountId string    `json:"virtual_account_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UserAccountsFilter struct {
	UserId  *string `json:"user_id"`
	Asset   *string `json:"asset"`
	Fiat    *string `json:"fiat"`
	Country *string `json:"country"`
}
