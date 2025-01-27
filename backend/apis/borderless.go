package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Deposit struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
	CreatedAt   time.Time   `json:"createdAt"`
	TxHash      *string     `json:"txHash"` // Use a pointer for nullable fields
	FeeAmount   string      `json:"feeAmount"`
}

type Source struct {
	Amount       string `json:"amount"`
	FiatCurrency string `json:"fiatCurrency"`
}

type Destination struct {
	Asset     string `json:"asset"`
	AccountID string `json:"accountId"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

// Borderless A struct for handling all borderless integrations
type Borderless struct {
	accessToken    string
	clientID       string
	clientSecret   string
	accountID      string
	BaseUrl        string
	Client         *http.Client
	Headers        map[string]interface{}
	Timeout        time.Duration
	IdempotencyKey uuid.UUID
}

// MakeRequest makes an HTTP request with retry logic and error handling
func (hc Borderless) MakeRequest(method, url string, data map[string]interface{}) (map[string]interface{}, error) {
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
	log.Printf("getting %s repsonse from %s with data: %v", method, url, data)
	if response.StatusCode >= 400 {
		fmt.Println("Response status code:", response.StatusCode)
		bodyBytes, _ := io.ReadAll(response.Body)
		var failedResponse map[string]interface{}
		if jsonErr := json.Unmarshal(bodyBytes, &failedResponse); jsonErr == nil {
			log.Printf("HTTP error occurred: %s", failedResponse)
			return failedResponse, errors.New("HTTP error")
		}
		fmt.Println("Failed response:", failedResponse)
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

func (hc Borderless) MakeDeposit(amount, asset, country, fiat string) (Deposit, error) {
	requestData := map[string]interface{}{
		"accountId":     hc.accountID,
		"amount":        amount,
		"asset":         asset,
		"country":       country,
		"fiat":          fiat,
		"paymentMethod": "wire",
	}
	// Make the request
	response, err := hc.MakeRequest(
		"POST",
		fmt.Sprintf("%s/deposits", hc.BaseUrl),
		requestData,
	)
	if err != nil {
		return Deposit{}, err
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		return Deposit{}, fmt.Errorf("failed to serialize deposit response: %w", err)
	}
	depositResponse := Deposit{}
	if err := json.Unmarshal(responseByte, &depositResponse); err != nil {
		return Deposit{}, fmt.Errorf("failed to parse deposit response: %w", err)
	}
	return depositResponse, nil

}

func (hc Borderless) GetTransaction(txId string) (map[string]interface{}, error) {
	// Make the request
	response, err := hc.MakeRequest(
		"GET",
		fmt.Sprintf("%s/transactions/%s", hc.BaseUrl, txId),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func NewBorderless() *Borderless {
	borderless := &Borderless{}
	clientID := os.Getenv("BORDERLESS_CLIENT_ID")
	if clientID == "" {
		log.Fatal("BORDERLESS_CLIENT_ID is required")
	}
	clientSecret := os.Getenv("BORDERLESS_CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("BORDERLESS_CLIENT_SECRET is required")
	}
	accountID := os.Getenv("BORDERLESS_ACCOUNT_ID")
	if accountID == "" {
		log.Fatal("BORDERLESS_ACCOUNT_ID is required")
	}
	borderless.accountID = accountID
	borderless.clientID = clientID
	borderless.clientSecret = clientSecret
	borderless.Headers = map[string]interface{}{
		"content-type": "application/json",
		"accept":       "application/json",
	}
	borderless.Timeout = 10 * time.Second
	borderless.BaseUrl = os.Getenv("BORDERLESS_BASE_URL")
	c := cache.New(60*time.Minute, 100*time.Minute)
	accessToken, found := c.Get("borderless_access_token")

	if !found {
		method := "POST"
		url := fmt.Sprintf("%s/auth/m2m/token", borderless.BaseUrl)
		borderless.Client = &http.Client{}
		data := map[string]interface{}{
			"clientId":     borderless.clientID,
			"clientSecret": borderless.clientSecret,
		}
		response, err := borderless.MakeRequest(
			method,
			url,
			data,
		)
		if err != nil {
			fmt.Println("Error making request:", err)
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling map:", err)
			log.Fatal(err)
		}

		// Unmarshal JSON to struct
		var token TokenResponse
		err = json.Unmarshal(jsonData, &token)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			log.Fatal(err)
		}
		accessToken = token.AccessToken
		c := cache.New(time.Duration(token.ExpiresIn-10)*time.Second, 10*time.Minute)
		c.Set("borderless_access_token", accessToken, time.Duration(token.ExpiresIn-10)*time.Second)
	}
	borderless.accessToken = accessToken.(string)
	borderless.Headers["idempotency-key"] = borderless.IdempotencyKey.String()
	return borderless
}
