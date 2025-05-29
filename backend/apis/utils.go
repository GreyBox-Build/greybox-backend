package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Rate struct {
	CurrencyName string `json:"currencyName"`
	Rate         string `json:"rate"`
}

// Define the structure for the data field
type RateData struct {
	UpdatedDate      time.Time       `json:"updated_date"`
	BaseCurrencyCode string          `json:"baseCurrencyCode"`
	Amount           int             `json:"amount"`
	Rates            map[string]Rate `json:"rates"`
}

// Define the main structure
type RateResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    RateData `json:"data"`
}

func GetExchangeRate(fiatCurrency, cryptoCurrency string, resultChan chan<- string, errChan chan<- error) {
	apiUrl := fmt.Sprintf("https://api.tatum.io/v3/tatum/rate/%s?basePair=%s", cryptoCurrency, fiatCurrency)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		errChan <- err
		return
	}
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()
	fmt.Println("status code", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		errChan <- err
		return
	}
	// Extract the exchange rate from the response
	exchangeRate, ok := data["value"].(string)
	if ok {
		resultChan <- exchangeRate
	} else {
		errChan <- fmt.Errorf("unexpected response format")
	}
}

func GetMobileMoneyExhangeRate(fiatCurrency string, resultChan chan<- string, errChan chan<- error) {
	apiUrl := fmt.Sprintf("https://sandbox.hurupay.com/v1/exchange/transfer_rate?from=USD&to=%s", fiatCurrency)
	// Create a new request
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		errChan <- err
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("HURUPAY_API_KEY")))
	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()
	var data RateResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		errChan <- err
		return
	}
	// Extract the exchange rate from the response
	_, ok := data.Data.Rates[fiatCurrency]
	if ok {
		resultChan <- data.Data.Rates[fiatCurrency].Rate
	} else {
		errChan <- fmt.Errorf("unexpected response format")
	}

}
