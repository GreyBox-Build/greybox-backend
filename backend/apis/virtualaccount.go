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

func CreateVirtualAccount(apiURL string, apiKey string, accountData serializers.VirtualAccount) error {
	// Convert struct to JSON
	jsonData, err := json.Marshal(accountData)
	if err != nil {
		return err
	}

	// Create HTTP client
	client := &http.Client{}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for 200 status code
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return errors.New(errorResponse.Message)
		}
		return errors.New(errorResponse.Message)
	}

	return nil
}

func FetchVirtualAccountDetail(id string, apiKey string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://api.tatum.io/v3/ledger/account/%s", id)

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
		// Return the error message
		return errData, fmt.Errorf(errorMessage)
	}
}
