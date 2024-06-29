package serializers

type OffRampForm struct {
	Amount         string `json:"amount"`
	AccountAddress string `json:"account_address"`
	Chain          string `json:"chain"`
}

type TransferXLM struct {
	To            string `json:"to"`
	Amount        string `json:"amount"`
	FromSecret    string `json:"fromSecret"`
	Initialize    bool   `json:"initialize"`
	Token         string `json:"token"`
	IssuerAccount string `json:"issuerAccount"`
}
