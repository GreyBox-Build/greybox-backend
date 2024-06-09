package serializers

type OffRampForm struct {
	Amount         string `json:"amount"`
	AccountAddress string `json:"account_address"`
	Chain          string `json:"chain"`
}
