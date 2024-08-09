package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

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
