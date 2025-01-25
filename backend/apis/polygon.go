package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	startTime := time.Now()
	var requestBody []byte
	var err error
	if data != nil {
		requestBody, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize request body: %w", err)
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range hc.Headers {
		req.Header.Set(key, value.(string))
	}
	log.Printf("Making %s request to %s with data: %v", method, url, data)
	var response *http.Response
	retries := 3
	for i := 0; i < retries; i++ {
		response, err = hc.Client.Do(req)
		if err == nil {
			break
		}
		if i < retries-1 {
			log.Printf("Retrying request to %s (attempt %d)", url, i+2)
			time.Sleep(2 * time.Second)
		}
	}
	// Handle errors from retries
	if err != nil {
		return nil, errors.New("an error occurred while making the request")
	}
	defer response.Body.Close()
	// Check for timeout or status code errors
	if response.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(response.Body)
		var failedResponse map[string]interface{}
		if jsonErr := json.Unmarshal(bodyBytes, &failedResponse); jsonErr == nil {
			log.Printf("HTTP error occurred: %s", failedResponse)
			return failedResponse, errors.New("HTTP error")
		}
		return failedResponse, errors.New("an unknown error occurred")
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}
	log.Printf("Successfully retrieved response in %v: %v", time.Since(startTime), responseData)
	return responseData, nil
}

func NewTatumPolygon() *TatumPolygon {
	baseUrl := os.Getenv("TATUM_BASE_URL")
	authKeys := os.Getenv("TATUM_API_KEY_TEST")
	return &TatumPolygon{
		BaseUrl: baseUrl,
		Client:  &http.Client{},
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"x-api-key":    authKeys,
			"accept":       "application/json",
			"timeout":      30,
		},
	}
}

func (hc *TatumPolygon) CreateWallet() (interface{}, error) {
	url := fmt.Sprintf("%s/polygon/wallet", hc.BaseUrl)
	response, err := hc.MakeRequest("POST", url, nil)
	if err != nil {
		return response, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := CreateWalletResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return response, nil
}

func (hc *TatumPolygon) GenerateAddress(xPub string, index uint) (interface{}, error) {
	url := fmt.Sprintf("%s/polygon/account/address/%s/%d", hc.BaseUrl, xPub, index)
	response, err := hc.MakeRequest("GET", url, nil)
	if err != nil {
		return response, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := GenerateAddressResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return response, nil
}

func (hc *TatumPolygon) GeneratePrivateKey(mnemonic string, index uint) (interface{}, error) {
	url := fmt.Sprintf("%s/polygon/wallet/priv", hc.BaseUrl)
	data := map[string]interface{}{
		"mnemonic": mnemonic,
		"index":    index,
	}
	response, err := hc.MakeRequest("POST", url, data)
	if err != nil {
		return response, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := GeneratePrivateKeyResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return response, nil
}

func (hc *TatumPolygon) GetAccountTransactions(address string, pageSize uint) (interface{}, error) {
	url := fmt.Sprintf("%s/polygon/account/transaction/%s?sort=%s&pageSize=%d", hc.BaseUrl, address, "DESC", pageSize)
	response, err := hc.MakeRequest("GET", url, nil)
	if err != nil {
		return response, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := make([]Transaction, 0)
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return response, nil
}

func (hc *TatumPolygon) GetTransaction(hash string) (interface{}, error) {
	url := fmt.Sprintf("%s/polygon/transaction/%s", hc.BaseUrl, hash)
	response, err := hc.MakeRequest("GET", url, nil)
	if err != nil {
		return response, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := Transaction{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return response, nil
}

// PerformTransaction uses tatum api to perform polygon transaction
func (hc *TatumPolygon) PerformTransaction(to, amount, privateKey string, currency CurrencyType) (interface{}, error) {
	url := fmt.Sprintf("%s/polygon/transaction", hc.BaseUrl)
	data := map[string]interface{}{
		"to":             to,
		"amount":         amount,
		"currency":       currency.String(),
		"fromPrivateKey": privateKey,
	}
	response, err := hc.MakeRequest("POST", url, data)
	if err != nil {
		return response, err
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %w", err)
	}
	responseData := TransactionHashResponse{}
	if err := json.Unmarshal(responseBytes, &responseData); err == nil {
		return responseData, nil
	}
	return response, nil
}
