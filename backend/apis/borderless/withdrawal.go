package borderless

import (
	"encoding/json"
	"fmt"
	"time"
)

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

type WithdrawalRequest struct {
	Fiat                 string `json:"fiat"`
	Country              string `json:"country"`
	Asset                string `json:"asset"`
	Amount               string `json:"amount"`
	AccountID            string `json:"accountId"`
	PaymentPurpose       string `json:"paymentPurpose"`
	PaymentInstructionID string `json:"paymentInstructionId"`
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
