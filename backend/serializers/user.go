package serializers

type User struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `gorm:"unique;not null" json:"email"`
	Password    string `json:"password"`
	Currency    string `json:"currency"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}

type LoginSerializer struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
