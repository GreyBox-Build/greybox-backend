package apis

import (
	"backend/serializers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type Result struct {
	Chain                  string `json:"chain"`
	Address                string `json:"address"`
	Balance                string `json:"balance"`
	TokenAddress           string `json:"tokenAddress"`
	LastUpdatedBlockNumber int64  `json:"lastUpdatedBlockNumber"`
	Type                   string `json:"type"`
	TokenId                string `json:"tokenId,omitempty"`
}

type Response struct {
	Result   []Result `json:"result"`
	PrevPage string   `json:"prevPage"`
	NextPage string   `json:"nextPage"`
}

type ManagedWallet struct {
	SignatureId string `json:"signatureId"`
}

type PrivateKeyKMS struct {
	PrivateKey string `json:"privateKey"`
}

type Wallet struct {
	Mnemonic string `json:"mnemonic"`
	Xpub     string `json:"xpub"`
	Testnet  bool   `json:"testnet"`
	Chain    string `json:"chain"`
}

type Address struct {
	Address string `json:"address"`
}

func GenerateCelloAddress(xpub string) (string, error) {
	apiURL := fmt.Sprintf("https://api.tatum.io/v3/celo/address/%s"+"/%d", xpub, 1)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	apiKey := os.Getenv("TATUM_API_KEY_TEST")

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	address := data["address"]

	return address, nil
}

func GenerateStellarAddress() (map[string]string, error) {
	apiUrl := "https://api.tatum.io/v3/xlm/account"
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	apiKey := os.Getenv("TATUM_API_KEY_TEST")
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GenerateCelloWallet() (string, string, error) {
	apiURL := "https://api.tatum.io/v3/celo/wallet"
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", "", err
	}
	apiKey := os.Getenv("TATUM_API_KEY_TEST")

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Parse JSON response
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return "", "", err
	}

	mnemonic := data["mnemonic"]
	xpub := data["xpub"]
	return mnemonic, xpub, nil
}

func GeneratePrivateKey(privData serializers.PrivGeneration) (string, error) {
	privURL := "https://api.tatum.io/v3/celo/wallet/priv"
	key := os.Getenv("TATUM_API_KEY_TEST")

	jsonData, err := json.Marshal(privData)
	if err != nil {
		return "", err
	}

	// Create HTTP client
	client := &http.Client{}

	// Create HTTP request
	req, err := http.NewRequest("POST", privURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", key)
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for 200 status code
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Message string        `json:"message"`
			Data    []interface{} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {

			return "", errors.New(errorResponse.Message)
		}

		return "", errors.New(errorResponse.Message)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	privkey := data["key"]
	if str, ok := privkey.(string); ok {
		return str, nil
	}
	return "", nil

}

func CreateVirtualAccount(apiURL string, apiKey string, accountData serializers.VirtualAccount) (string, error) {
	// Convert struct to JSON
	jsonData, err := json.Marshal(accountData)
	if err != nil {
		return "", err
	}

	// Create HTTP client
	client := &http.Client{}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for 200 status code
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Message string        `json:"message"`
			Data    []interface{} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return "", errors.New(errorResponse.Message)
		}

		return "", errors.New(errorResponse.Message)
	}
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	customerId, ok := data["id"].(string)
	if ok {
		return customerId, nil
	}

	return "", nil
}

func FetchAccountBalance(id string, apiKey string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://api.tatum.io/v3/ledger/account/%s"+"/balance", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("x-api-key", apiKey)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode == http.StatusOK {
		// Decode the JSON response
		var data map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		// Decode the error JSON response
		var errData map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errData)
		if err != nil {
			return nil, err
		}

		// Extract the error message from the JSON response
		errorMessage, ok := errData["message"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to extract error message from json response")
		}
		errorData, ok := errData["data"].([]interface{})
		if ok {
			fmt.Println("error data: ", errorData)
		}
		// Return the error message
		return errData, fmt.Errorf(errorMessage)
	}
}

func ActivateVirtualAccount(id string, apiKey string) error {

	url := fmt.Sprintf("https://api.tatum.io/v3/ledger/account/%s"+"/activate", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("x-api-key", apiKey)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode == http.StatusOK {

		return nil
	} else if resp.StatusCode == http.StatusBadRequest {
		return errors.New("request validation failed")
	} else if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("subscription not active anymore")
	} else if resp.StatusCode == http.StatusForbidden {
		return errors.New("no such account")
	} else if resp.StatusCode == http.StatusInternalServerError {
		return errors.New("internal server error")
	}

	return nil
}

func CreateDepositWallet(accountId string) (string, string, error) {

	apiUrl := fmt.Sprintf("https://api.tatum.io/v3/offchain/account/%s/address", accountId)
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return "", "", err
	}
	apiKey := os.Getenv("TATUM_API_KEY_TEST")
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", err
	}
	address := data["address"]
	message := data["message"]
	str, _ := address.(string)
	mes, _ := message.(string)
	return str, mes, nil

}

func FetchWalletBalance(address, chain string, pageSize int32) (float32, error) {
	tokenType := "fungible"
	apiUrl := fmt.Sprintf("https://api.tatum.io/v4/data/balances?chain=%s&addresses=%s&excludeMetadata=%t&tokenTypes=%s&pageSize=%d", chain, address, true, tokenType, 10)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return 0, err
	}
	apiKey := os.Getenv("TATUM_API_KEY_TEST")
	req.Header.Add("x-api-key", apiKey)
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	var balance float32
	for _, result := range data.Result {
		balance1, err := strconv.ParseFloat(result.Balance, 32)
		if err != nil {
			return 0, err
		}
		balance = float32(balance1)

	}
	return balance, nil
}

func GenerateXlmAccount() (map[string]string, int, error) {
	apiUrl := "https://api.tatum.io/v3/xlm/account"

	client := &http.Client{}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, 500, err
	}
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		var data map[string]string
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, 500, err
		}
		return data, 200, nil
	case 400:
		return nil, 400, errors.New("bad Request. Validation Error")
	case 401:
		return nil, 401, errors.New("unauthorized")
	case 500:
		return nil, 500, errors.New("internal Server Error")
	default:
		return nil, 403, errors.New("unable to communicate with blockchain")
	}

}

func FetchAccountBalanceXLM(address string) (float32, error) {
	apiUrl := fmt.Sprintf("https://api.tatum.io/v3/xlm/account/%s", address)

	client := &http.Client{}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	respData := serializers.Account{}
	fmt.Println("code: ", resp.StatusCode)
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return 0, err
	}
	switch resp.StatusCode {
	case 400:
		return 0, errors.New("bad request")
	case 500:
		return 0, errors.New("internal server error")
	case 401:
		return 0, errors.New("subscription might not be active again")
	case 403:
		return 0, errors.New("unable to communicate with blockchain")

	case 404:
		return 0, nil
	default:
		amount := 0.0
		for _, balance := range respData.Balances {
			fmt.Println("balance: ", balance.Balance, balance.AssetType)
			if balance.AssetType != "native" {
				a, _ := strconv.ParseFloat(balance.Balance, 32)
				amount += a
			}

		}
		return float32(amount), nil
	}

}

func FetchAccountBalanceCUSD(address string) (float32, error) {
	apiUrl := fmt.Sprintf("https://api.tatum.io/v3/celo/account/balance/%s", address)

	client := &http.Client{}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	respData := map[string]string{}
	switch resp.StatusCode {
	case 400:
		return 0, errors.New("bad request")
	case 500:
		return 0, errors.New("internal server error")
	case 401:
		return 0, errors.New("subscription might not be active again")
	case 403:
		return 0, errors.New("unable to communicate with blockchain")

	case 404:
		return 0, nil
	default:
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return 0, err
		}
		amount, _ := strconv.ParseFloat(respData["cUsd"], 32)
		return float32(amount), nil
	}

}

func StorePrivateKeyManagedWallet(asset string) (ManagedWallet, error) {
	cmd := exec.Command("tatum-kms", "storemanagedprivatekey", asset)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ManagedWallet{}, fmt.Errorf("failed to execute command: %w", err)
	}
	stringOutput := string(output)
	var wallet ManagedWallet
	if err := json.Unmarshal([]byte(stringOutput), &wallet); err != nil {
		return ManagedWallet{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	privateKeyCMD := exec.Command("tatum-kms", "storemanagedprivatekey", asset)
	_, err = privateKeyCMD.CombinedOutput()
	if err != nil {
		return ManagedWallet{}, fmt.Errorf("failed to execute command: %w", err)
	}
	return wallet, nil
}

func GetPrivateKeyManagedWallet(signatureId string, index uint64) (PrivateKeyKMS, error) {
	stringIndex := strconv.FormatUint(index, 10)

	cmd := exec.Command("tatum-kms", "getprivatekey", signatureId, stringIndex)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return PrivateKeyKMS{}, fmt.Errorf("failed to execute command: %w", err)
	}
	stringOutput := string(output)
	var privKey PrivateKeyKMS
	if err := json.Unmarshal([]byte(stringOutput), &privKey); err != nil {
		return PrivateKeyKMS{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return privKey, nil
}

func GetManagedWallet(signatureId string) (Wallet, error) {
	cmd := exec.Command("tatum-kms", "getmanagedwallet", signatureId)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to execute command: %w", err)
	}
	stringOutput := string(output)
	var wallet Wallet
	if err := json.Unmarshal([]byte(stringOutput), &wallet); err != nil {
		return Wallet{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return wallet, nil
}

func GetManagedWalletAddress(signatureId string, index uint64) (Address, error) {
	stringIndex := strconv.FormatUint(index, 10)
	cmd := exec.Command("tatum-kms", "getaddress", signatureId, stringIndex)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Address{}, fmt.Errorf("failed to execute command: %w", err)
	}
	stringOutput := string(output)
	var address Address
	if err := json.Unmarshal([]byte(stringOutput), &address); err != nil {
		return Address{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return address, nil

}
