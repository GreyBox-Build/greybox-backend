package borderless

import (
	"backend/state"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

type DepositOrWithdrawalOption struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Method string `json:"method"`
	Fiat   string `json:"fiat"`
	Asset  string `json:"asset"`
}

func NewBorderless() *Borderless {
	borderless := &Borderless{}
	clientID := state.AppConfig.BorderlessClientId

	clientSecret := state.AppConfig.BorderlessClientSecret

	accountID := state.AppConfig.BorderlessAccountId

	borderless.accountID = accountID
	borderless.clientID = clientID
	borderless.clientSecret = clientSecret
	borderless.Headers = map[string]interface{}{
		"content-type": "application/json",
		"accept":       "application/json",
	}
	borderless.Timeout = 10 * time.Second
	borderless.BaseUrl = state.AppConfig.BorderlessBaseUrl

	// Check for token in the cache
	accessToken, found := state.BorderlessCache.Get("borderless_access_token")
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
			log.Println("Error making request:", err)
			log.Fatal(err)
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshalling map:", err)
			log.Fatal(err)
		}

		// Unmarshal JSON to struct
		var token TokenResponse
		err = json.Unmarshal(jsonData, &token)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			log.Fatal(err)
		}

		accessToken = token.AccessToken
		state.BorderlessCache.Set(
			"borderless_access_token",
			accessToken,
			time.Duration(token.ExpiresIn-10)*time.Second,
		)
	}

	borderless.accessToken = accessToken.(string)
	borderless.Headers["Authorization"] = fmt.Sprintf("Bearer %s", borderless.accessToken)
	return borderless
}

func (hc Borderless) GetAvailableCountries(entity string) ([]string, error) {
	response, err := hc.MakeRequest("GET", fmt.Sprintf("%s/%s/countries", hc.BaseUrl, entity), nil)

	if err != nil {
		return nil, err
	}

	// Check if the response is of the expected type
	fmt.Printf("Countries Response: %v", response)
	responseBytes, err := json.Marshal(response)
	fmt.Printf("Response Bytes: %v", responseBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var countries []string
	err = json.Unmarshal(responseBytes, &countries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal countries: %w", err)
	}

	return countries, nil
}

func (hc Borderless) GetDepositOrWithdrawalOption(
	entity string,
	countryCode string,
	fiat string,
	asset string,
) (*DepositOrWithdrawalOption, error) {
	response, err := hc.MakeRequest("GET", fmt.Sprintf("%s/%s/options?country=%s", hc.BaseUrl, entity, countryCode), nil)

	if err != nil {
		return nil, err
	}

	// Check if the response is of the expected type
	fmt.Printf("Country Options Response: %v", response)
	responseBytes, err := json.Marshal(response)
	fmt.Printf("Response Bytes: %v", responseBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	var options []DepositOrWithdrawalOption
	err = json.Unmarshal(responseBytes, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s options: %w", entity, err)
	}

	for _, option := range options {
		if strings.EqualFold(option.Method, "mobile money") || strings.EqualFold(option.Method, "mobilemoney") {
			if strings.EqualFold(option.Fiat, fiat) && strings.EqualFold(option.Asset, asset) {
				if option.ID == "" {
					return nil, errors.New(fmt.Sprintf("mobile money %s option found without ID", entity))
				}
				return &option, nil
			}
		}
	}

	return nil, fmt.Errorf("No Mobile Money %s option found for %s: %s and %s", entity, countryCode, fiat, asset)
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
