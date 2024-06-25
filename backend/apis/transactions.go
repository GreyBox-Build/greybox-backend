package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	
)

type TransactionRequest struct {
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	To             string `json:"to"`
	FeeCurrency    string `json:"feeCurrency"`
	FromPrivateKey string `json:"fromPrivateKey"` // Omit this field from JSON
}

type TransactionRequestV2 struct {
	FeeCurrency      string `json:"feeCurrency"`
	Chain            string `json:"chain"`
	CustodialAddress string `json:"custodialAddress"`
	Recipient        string `json:"recipient"`
	ContractType     uint32 `json:"contractType"`
	TokenAddress     string `json:"tokenAddress"`
	Amount           string `json:"amount"`
	FromPrivateKey   string `json:"fromPrivateKey"`
	Fee              struct {
		GasPrice string `json:"gasPrice"`
		GasLimit string `json:"gasLimit"`
	} `json:"fee"`
}

type ErrorResponse struct {
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       []Data `json:"data"`
}

type Data struct {
	Target      Target      `json:"target"`
	Value       int         `json:"value"`
	Property    string      `json:"property"`
	Constraints Constraints `json:"constraints"`
}

type Target struct {
	Property int `json:"property"`
}

type Constraints struct {
	Min string `json:"min"`
}

func GetUserTransactions(chain, walletAddress, category string, pageSize uint64) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.tatum.io/v4/data/transactions?chain=%s&addresses=%s&transactionSubtype=%s&pageSize=%d", chain, walletAddress, category, pageSize)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get user transactions")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetTransactionByHash(chain, hash string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.tatum.io/v4/data/transactions/hash?hash=%s&chain=%s", hash, chain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get user transactions")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func PerformTransactionCelo(amount, accountAddress, privKey, gasPrice string, gasFee float64) (string, int, error) {
	url := "https://api.tatum.io/v3/celo/transaction"
	client := &http.Client{}

	

	newData := TransactionRequest{
		Amount:         amount,
		Currency:       "CUSD",
		To:             accountAddress,
		FeeCurrency:    "CELO",
		FromPrivateKey: privKey,
	}

	// Convert the struct to JSON format
	jsonData, err := json.Marshal(&newData)
	if err != nil {
		return "", 500, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 500, err
	}

	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", 500, err
	}
	defer resp.Body.Close()
	var errMsg string

	switch resp.StatusCode {
	case 403:
		errMsg = "failed to perform transaction. most likely insufficient funds"

	case 400:
		errMsg = "failed to perform transaction. validation error"

	case 500:
		errMsg = "server error from third party application"
	case 401:
		errMsg = "subscription not active"
	default:
		// Handle any other status codes if needed
	}

	if errMsg != "" {
		return "", resp.StatusCode, errors.New(errMsg)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 500, err
	}
	var result map[string]string
	if err = json.Unmarshal(body, &result); err != nil {
		return "", 500, err
	}
	return result["txId"], resp.StatusCode, nil
}

func CalculateEstimatedFeeCelo(amount, to, from string) (map[string]interface{}, error) {
	apiUrl := "https://api.tatum.io/v3/celo/gas"
	client := &http.Client{}
	jsonData, err := json.Marshal(map[string]string{
		"amount": amount,
		"to":     to,
		"from":   from,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to calculate estimated fee")
	}
	var result map[string]interface{}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
