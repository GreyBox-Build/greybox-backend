package serializers

import (
	"time"
)

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
	FromAccount   string `json:"fromAccount"`
}

type Transaction struct {
	Chain              string `json:"chain"`
	Hash               string `json:"hash"`
	Address            string `json:"address"`
	BlockNumber        int    `json:"blockNumber"`
	TransactionIndex   int    `json:"transactionIndex"`
	TransactionType    string `json:"transactionType"`
	TransactionSubtype string `json:"transactionSubtype"`
	Amount             string `json:"amount"`
	Timestamp          int64  `json:"timestamp"`
	TokenAddress       string `json:"tokenAddress,omitempty"`
	CounterAddress     string `json:"counterAddress"`
	TokenId            string `json:"tokenId,omitempty"`
}

type Result struct {
	Result   []Transaction `json:"result"`
	PrevPage string        `json:"prevPage"`
	NextPage string        `json:"nextPage"`
}

type TransactionXLM struct {
	ID                    string    `json:"id"`
	PagingToken           string    `json:"paging_token"`
	Successful            bool      `json:"successful"`
	Hash                  string    `json:"hash"`
	Ledger                int       `json:"ledger"`
	CreatedAt             time.Time `json:"created_at"`
	SourceAccount         string    `json:"source_account"`
	SourceAccountSequence string    `json:"source_account_sequence"`
	FeePaid               int       `json:"fee_paid"`
	FeeCharged            string    `json:"fee_charged"`
	MaxFee                string    `json:"max_fee"`
	OperationCount        int       `json:"operation_count"`
	EnvelopeXDR           string    `json:"envelope_xdr"`
	ResultXDR             string    `json:"result_xdr"`
	ResultMetaXDR         string    `json:"result_meta_xdr"`
	FeeMetaXDR            string    `json:"fee_meta_xdr"`
	Memo                  string    `json:"memo"`
	MemoType              string    `json:"memo_type"`
	Signatures            []string  `json:"signatures"`
}

type SignUrl struct {
	Url string `json:"url"`
}

type Bank struct {
	CountryCode   string `json:"CountryCode"`
	BankName      string `json:"BankName"`
	AccountNumber string `json:"AccountNumber"`
	AccountName   string `json:"AccountName"`
}

type BanksStruct struct {
	Banks []Bank `json:"Banks"`
}
