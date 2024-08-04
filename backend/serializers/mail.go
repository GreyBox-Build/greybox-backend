package serializers

type AdminOnRampSerializer struct {
	Name          string
	BankName      string
	AccountNumber string
	AccountName   string
	Currency      string
	Ref           string
	Amount        string
}

type AdminOffRampSerializer struct {
	Name          string
	BankName      string
	AccountNumber string
	Currency      string
	Amount        string
	Ref           string
}

type UserOffRampMail struct {
	Amount        string
	Name          string
	Ref           string
	BankName      string
	AccountNumber string
	AccountName   string
	Currency      string
}
