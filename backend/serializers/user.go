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
