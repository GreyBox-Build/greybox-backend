package borderless

import (
	"backend/state"
	"encoding/json"
	"fmt"
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

type PaymentResponse struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	PaymentMethod string              `json:"paymentMethod"`
	Currency      string              `json:"currency"`
	Country       string              `json:"country"`
	Deleted       bool                `json:"deleted"`
	Details       BankDetailsResponse `json:"details"`
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

func (hc Borderless) MakePaymentInstruction(payment Payment) (PaymentResponse, error) {
	data, err := StructToMap(payment)
	if err != nil {
		return PaymentResponse{}, err
	}
	identity := state.AppConfig.BorderlessBusinessIdentity
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

func NewPayment(
	country,
	currency,
	name,
	paymentMethod,
	accountHolderName,
	bankAccountNumber,
	bankAccountType,
	bankName,
	street1,
	city,
	addressCountry,
	postalCode string,
	bankRoutingNumber,
	iban,
	swiftOrBicNumber,
	street2,
	state string) Payment {
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
