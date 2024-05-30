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
)

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

func GeneratePrivateKey(apiURL string, apiKey string, privData serializers.PrivGeneration) (string, error) {
	// Convert struct to JSON
	jsonData, err := json.Marshal(privData)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	key := data["key"]
	if str, ok := key.(string); ok {
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
