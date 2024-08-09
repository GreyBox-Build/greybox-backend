package serializers

type VirtualAccount struct {
	Xpub               string            `json:"xpub"`
	Currency           string            `json:"currency"`
	Customer           map[string]string `json:"customer"`
	Compliant          bool              `json:"compliant"`
	AccountCode        string            `json:"accountCode"`
	AccountingCurrency string            `json:"accountingCurrency"`
	AccountNumber      string            `json:"accountNumber"`
}

type PrivGeneration struct {
	Index    uint   `json:"index"`
	Mnemonic string `json:"mnemonic"`
}

type MasterWalletForm struct {
	Asset string `json:"asset"`
}
