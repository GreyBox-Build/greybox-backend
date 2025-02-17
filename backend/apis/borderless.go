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

type PaymentAddress struct {
	Street1    string
	Street2    string
	City       string
	State      string
	Country    string
	PostalCode string
}

type BankDetails struct {
	AccountHolderName      string
	BankAccountNumber      string
	BankAccountNumberLast4 string
	BankAccountType        string // "Checking" or "Savings"
	BankName               string
	BankRoutingNumber      string
	IBAN                   string
	SwiftOrBicNumber       string
	Address                PaymentAddress
}

type Payment struct {
	Country       string
	Currency      string
	Name          string
	PaymentMethod string // "ACH"
	Details       BankDetails
}

var borderlessCache = cache.New(60*time.Minute, 100*time.Minute)

type Deposit struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
	CreatedAt   time.Time   `json:"createdAt"`
	TxHash      *[]string   `json:"txHash"`
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
type AddressResponse struct {
	ID         string `json:"id"`
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type BankDetailsResponse struct {
	BankName               string          `json:"bankName"`
	BankAccountNumberLast4 int             `json:"bankAccountNumberLast4"`
	BankRoutingNumber      string          `json:"bankRoutingNumber"`
	AccountHolderName      string          `json:"accountHolderName"`
	BankAccountType        string          `json:"bankAccountType"`
	Address                AddressResponse `json:"address"`
}

type PaymentResponse struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	PaymentMethod string              `json:"paymentMethod"`
	Currency      string              `json:"currency"`
	Country       string              `json:"country"`
	Deleted       bool                `json:"deleted"`
	Details       BankDetailsResponse `json:"details"`
}

type WithdrawalSource struct {
	FiatCurrency string `json:"fiatCurrency"`
	Amount       string `json:"amount"`
}

type WithdrawalDestination struct {
	Currency  string `json:"currency"`
	AccountID string `json:"accountId"`
}

type WithdrawalResponse struct {
	ID          string                `json:"id"`
	Type        string                `json:"type"`
	Status      string                `json:"status"`
	Source      WithdrawalSource      `json:"source"`
	Destination WithdrawalDestination `json:"destination"`
	CreatedAt   time.Time             `json:"createdAt"`
	TxHash      string                `json:"txHash"`
}

// Borderless A struct for handling all borderless integrations
type Borderless struct {
	accessToken  string
	clientID     string
	clientSecret string
	accountID    string
	BaseUrl      string
	Client       *http.Client
	Headers      map[string]interface{}
	Timeout      time.Duration
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
	hc.Client = &http.Client{}
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
		"paymentMethod": "Wire",
	}
	idempotencyKey := uuid.New()
	fmt.Println("Idempotency key:", idempotencyKey)
	hc.Headers["idempotency-key"] = idempotencyKey.String()
	// Make the request
	response, err := hc.MakeRequest(
		"POST",
		fmt.Sprintf("%s/deposits", hc.BaseUrl),
		requestData,
	)
	fmt.Println("Response:", response)
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

func (hc Borderless) MakePaymentInstruction(payment Payment) (PaymentResponse, error) {
	data, err := StructToMap(payment)
	if err != nil {
		return PaymentResponse{}, err
	}
	identity := os.Getenv("BORDERLESS_BUSINESS_IDENTITY")
	response, err := hc.MakeRequest(
		"POST", fmt.Sprintf("%s/identities/%s/payment-instructions", hc.BaseUrl, identity), data)
	if err != nil {
		return PaymentResponse{}, err
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		return PaymentResponse{}, fmt.Errorf("failed to serialize payment instruction response: %w", err)
	}
	paymentResponse := PaymentResponse{}
	if err := json.Unmarshal(responseByte, &paymentResponse); err != nil {
		return PaymentResponse{}, fmt.Errorf("failed to parse payment instruction response: %w", err)
	}
	return paymentResponse, nil

}

func (hc Borderless) MakeWithdrawal(request WithdrawalRequest) (WithdrawalResponse, error) {
	data, err := StructToMap(request)
	if err != nil {
		return WithdrawalResponse{}, err
	}
	response, err := hc.MakeRequest(
		"POST", fmt.Sprintf("%s/withdrawals", hc.BaseUrl), data)
	if err != nil {
		return WithdrawalResponse{}, err
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		return WithdrawalResponse{}, fmt.Errorf("failed to serialize withdrawal response: %w", err)
	}
	paymentResponse := WithdrawalResponse{}
	if err := json.Unmarshal(responseByte, &paymentResponse); err != nil {
		return WithdrawalResponse{}, fmt.Errorf("failed to parse withdrawal response: %w", err)
	}
	return paymentResponse, nil
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

	// Check for token in the cache
	accessToken, found := borderlessCache.Get("borderless_access_token")
	if !found {
		method := "POST"
		url := fmt.Sprintf("%s/auth/m2m/token", borderless.BaseUrl)
		borderless.Client = &http.Client{}
		data := map[string]interface{}{
			"clientId":     borderless.clientID,
			"clientSecret": borderless.clientSecret,
		}
		response, err := borderless.MakeRequest(method, url, data)
		if err != nil {
			fmt.Println("Error making request:", err)
			log.Fatal(err)
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
		borderlessCache.Set(
			"borderless_access_token",
			accessToken,
			time.Duration(token.ExpiresIn-10)*time.Second,
		)
	}

	borderless.accessToken = accessToken.(string)
	borderless.Headers["Authorization"] = fmt.Sprintf("Bearer %s", borderless.accessToken)
	return borderless
}

func NewPayment(country, currency, name, paymentMethod, accountHolderName, bankAccountNumber, bankAccountType, bankName, street1, city, addressCountry, postalCode string, bankRoutingNumber, iban, swiftOrBicNumber, street2, state string) Payment {
	if bankAccountType != "Checking" && bankAccountType != "Savings" {
		panic("Invalid bank account type. Must be 'Checking' or 'Savings'")
	}

	var bankAccountNumberLast4 string
	if len(bankAccountNumber) >= 4 {
		bankAccountNumberLast4 = bankAccountNumber[len(bankAccountNumber)-4:]
	}

	if street1 != "" {
		street2 = ""
	}

	return Payment{
		Country:       country,
		Currency:      currency,
		Name:          name,
		PaymentMethod: paymentMethod,
		Details: BankDetails{
			AccountHolderName:      accountHolderName,
			BankAccountNumber:      bankAccountNumber,
			BankAccountNumberLast4: bankAccountNumberLast4,
			BankAccountType:        bankAccountType,
			BankName:               bankName,
			BankRoutingNumber:      bankRoutingNumber,
			IBAN:                   iban,
			SwiftOrBicNumber:       swiftOrBicNumber,
			Address: PaymentAddress{
				Street1:    street1,
				Street2:    street2,
				City:       city,
				State:      state,
				Country:    addressCountry,
				PostalCode: postalCode,
			},
		},
	}
}

func StructToMap[T any](s T) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type WithdrawalRequest struct {
	Fiat                 string `json:"fiat"`
	Country              string `json:"country"`
	Asset                string `json:"asset"`
	Amount               string `json:"amount"`
	AccountID            string `json:"accountId"`
	PaymentPurpose       string `json:"paymentPurpose"`
	PaymentInstructionID string `json:"paymentInstructionId"`
}

func NewWithdrawalRequest(fiat, country, asset, amount, accountID, paymentPurpose, paymentInstructionID string) WithdrawalRequest {
	return WithdrawalRequest{
		Fiat:                 fiat,
		Country:              country,
		Asset:                asset,
		Amount:               amount,
		AccountID:            accountID,
		PaymentPurpose:       paymentPurpose,
		PaymentInstructionID: paymentInstructionID,
	}
}
