package apis

import (
	"backend/state"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type CurrencyType int

const (
	// MATIC Enum values for CurrencyType
	MATIC CurrencyType = iota
	USDC_MATIC
	USDT_MATIC
	INTENT
)

// String method for CurrencyType to get the string representation
func (c CurrencyType) String() string {
	switch c {
	case MATIC:
		return "MATIC"
	case USDC_MATIC:
		return "USDC_MATIC"
	case USDT_MATIC:
		return "USDT_MATIC"
	case INTENT:
		return "INTENT"
	default:
		return "UNKNOWN"
	}
}

// ParseCurrencyType converts a string to a CurrencyType
func ParseCurrencyType(value string) (CurrencyType, error) {
	switch value {
	case "MATIC":
		return MATIC, nil
	case "USDC_MATIC":
		return USDC_MATIC, nil
	case "USDT_MATIC":
		return USDT_MATIC, nil
	case "INTENT":
		return INTENT, nil
	default:
		return 0, fmt.Errorf("invalid CurrencyType: %s", value)
	}
}

type TransactionHashResponse struct {
	TxId string `json:"txId"`
}

type Transaction struct {
	BlockHash         string `json:"blockHash"`
	Status            bool   `json:"status"`
	BlockNumber       int    `json:"blockNumber"`
	From              string `json:"from"`
	Gas               int    `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	TransactionHash   string `json:"transactionHash"`
	Input             string `json:"input"`
	Nonce             int    `json:"nonce"`
	To                string `json:"to"`
	TransactionIndex  int    `json:"transactionIndex"`
	Value             string `json:"value"`
	GasUsed           int    `json:"gasUsed"`
	CumulativeGasUsed int    `json:"cumulativeGasUsed"`
	ContractAddress   string `json:"contractAddress"`
	Logs              []Log  `json:"logs"`
}

type Log struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	LogIndex         int      `json:"logIndex"`
	BlockNumber      int      `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	TransactionIndex int      `json:"transactionIndex"`
	TransactionHash  string   `json:"transactionHash"`
}

type TatumPolygon struct {
	BaseUrl string
	Client  *http.Client
	Headers map[string]interface{}
}

type CreateWalletResponse struct {
	Xpub     string `json:"xpub"`
	Mnemonic string `json:"mnemonic"`
}

type GenerateAddressResponse struct {
	Address string `json:"address"`
}

type GeneratePrivateKeyResponse struct {
	Key string `json:"key"`
}

// MakeRequest makes an HTTP request with retry logic and error handling
func (hc *TatumPolygon) MakeRequest(method, url string, data map[string]interface{}) (map[string]interface{}, error) {
	start := time.Now()

	// Marshal request body if data is present
	var requestBody []byte
	if data != nil {
		var err error
		requestBody, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize request body: %w", err)
		}
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range hc.Headers {
		if strVal, ok := value.(string); ok {
			req.Header.Set(key, strVal)
		}
	}

	log.Printf("Making %s request to %s with data: %v", method, url, data)

	// Attempt with retries
	const maxRetries = 3
	var resp *http.Response
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = hc.Client.Do(req)
		if err == nil {
			break
		}
		if attempt < maxRetries {
			log.Printf("Retrying request to %s (attempt %d)", url, attempt+1)
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-2xx status codes
	if resp.StatusCode >= 400 {
		var errorResponse map[string]interface{}
		if jsonErr := json.Unmarshal(body, &errorResponse); jsonErr == nil {
			log.Printf("HTTP %d error response: %v", resp.StatusCode, errorResponse)
			msg := "HTTP error"
			if m, ok := errorResponse["message"].(string); ok {
				msg = m
			}
			return errorResponse, fmt.Errorf("HTTP error: %s", msg)
		}
		log.Printf("Unparseable HTTP %d error body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("HTTP error: %s", string(body))
	}

	// Parse success response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	log.Printf("Successfully completed %s request to %s in %v", method, url, time.Since(start))
	return result, nil
}

func NewTatumPolygon() *TatumPolygon {
	baseUrl := state.AppConfig.TatumBaseUrl
	authKeys := state.AppConfig.TatumTestApiKey
	return &TatumPolygon{
		BaseUrl: baseUrl,
		Client:  &http.Client{},
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"x-api-key":    authKeys,
			"accept":       "application/json",
		},
	}
}

func (hc *TatumPolygon) CreateWallet() (CreateWalletResponse, error) {
	url := fmt.Sprintf("%s/polygon/wallet", hc.BaseUrl)
	response, err := hc.MakeRequest("GET", url, nil)
	if err != nil {
		return CreateWalletResponse{}, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return CreateWalletResponse{}, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := CreateWalletResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return CreateWalletResponse{}, nil
}

func (hc *TatumPolygon) GenerateAddress(xPub string, index uint) (GenerateAddressResponse, error) {
	url := fmt.Sprintf("%s/polygon/address/%s/%d", hc.BaseUrl, xPub, index)
	response, err := hc.MakeRequest("GET", url, nil)
	if err != nil {
		return GenerateAddressResponse{}, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return GenerateAddressResponse{}, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := GenerateAddressResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return GenerateAddressResponse{}, nil
}

func (hc *TatumPolygon) GeneratePrivateKey(mnemonic string, index uint) (GeneratePrivateKeyResponse, error) {
	url := fmt.Sprintf("%s/polygon/wallet/priv", hc.BaseUrl)
	data := map[string]interface{}{
		"mnemonic": mnemonic,
		"index":    index,
	}
	response, err := hc.MakeRequest("POST", url, data)
	if err != nil {
		return GeneratePrivateKeyResponse{}, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return GeneratePrivateKeyResponse{}, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := GeneratePrivateKeyResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return GeneratePrivateKeyResponse{}, nil
}

func (hc *TatumPolygon) GetAccountTransactions(address string, pageSize uint) ([]Transaction, error) {
	url := fmt.Sprintf("%s/polygon/account/transaction/%s?sort=%s&pageSize=%d", hc.BaseUrl, address, "DESC", pageSize)
	response, err := hc.MakeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := make([]Transaction, 0)
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return nil, nil
}

func (hc *TatumPolygon) GetTransaction(hash string) (Transaction, error) {
	url := fmt.Sprintf("%s/polygon/transaction/%s", hc.BaseUrl, hash)
	response, err := hc.MakeRequest("GET", url, nil)
	if err != nil {
		return Transaction{}, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := Transaction{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return Transaction{}, nil
}

// PerformTransaction uses tatum api to perform polygon transaction
func (hc *TatumPolygon) PerformTransaction(to, amount, privateKey string, currency CurrencyType) (TransactionHashResponse, error) {
	url := fmt.Sprintf("%s/polygon/transaction", hc.BaseUrl)
	data := map[string]interface{}{
		"to":             to,
		"amount":         amount,
		"currency":       currency.String(),
		"fromPrivateKey": privateKey,
	}
	response, err := hc.MakeRequest("POST", url, data)
	if err != nil {
		return TransactionHashResponse{}, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return TransactionHashResponse{}, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := TransactionHashResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return TransactionHashResponse{}, nil
}
