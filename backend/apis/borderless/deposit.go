package borderless

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Source struct {
	Amount       string `json:"amount"`
	FiatCurrency string `json:"fiatCurrency"`
}

type Destination struct {
	Asset     string `json:"asset"`
	AccountID string `json:"accountId"`
}

type Deposit struct {
	ID                    string      `json:"id"`
	SourceAccountId       string      `json:"sourceAccountId"`
	DestinationAccountId  interface{} `json:"destinationAccountId"` // Can be null
	Type                  string      `json:"type"`
	Status                string      `json:"status"`
	Source                Source      `json:"source"`
	Destination           Destination `json:"destination"`
	ProviderTransactionId interface{} `json:"providerTransactionId"` // Can be null
	Instructions          interface{} `json:"instructions"`          // Can be null
	FiatCurrency          string      `json:"fiatCurrency"`
	CreatedAt             time.Time   `json:"createdAt"`
	TxHash                *[]string   `json:"txHash"`
	FeeAmount             string      `json:"feeAmount"`
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

	log.Println("Idempotency Key: ", idempotencyKey)
	hc.Headers["idempotency-key"] = idempotencyKey.String()

	// Make the request
	response, err := hc.MakeRequest(
		"POST",
		fmt.Sprintf("%s/deposits", hc.BaseUrl),
		requestData,
	)

	log.Println("Response: ", response)

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

func (hc Borderless) MobileMoneyDeposit(
	accountId string,
	fiat string,
	country string,
	asset string,
	amount string,
	paymentMethod string,
) (Deposit, error) {

	// by default use greybox account ID but if an account ID is provided, use that instead
	innerAccountId := hc.accountID
	if accountId != "" {
		innerAccountId = accountId
	}

	requestData := map[string]interface{}{
		"accountId":     innerAccountId,
		"amount":        amount,
		"asset":         asset,
		"country":       country,
		"fiat":          fiat,
		"paymentMethod": paymentMethod,
	}

	// Make the request
	response, err := hc.MakeRequest(
		"POST",
		fmt.Sprintf("%s/deposits", hc.BaseUrl),
		requestData,
	)

	log.Println("Response: ", response)

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
