package serializers

type Chain struct {
	Celo    string `json:"celo"`
	Stellar string `json:"stellar"`
}

const (
	ChainCelo    = "CELO"
	ChainStellar = "XLM"
)

var Chains = Chain{
	Celo:    ChainCelo,
	Stellar: ChainStellar,
}

type User struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `gorm:"unique;not null" json:"email"`
	Password    string `json:"password"`
	Currency    string `json:"currency"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Chain       string `json:"chain"`
}

type LoginSerializer struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Data struct {
	Chain          string `json:"chain"`
	Name           string `json:"name"`
	ImageURL       string `json:"image_url"`
	AvailableToken string `json:"availableToken"`
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
	Key       string `json:"key"`
	UserEmail string `json:"user_email"`
}

type NetworkData struct {
	CountryCode  string   `json:"countryCode"`
	CountryName  string   `json:"countryName"`
	CurrencyCode string   `json:"currencyCode"`
	Networks     []string `json:"networks"`
}
