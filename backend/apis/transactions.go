package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetUserTransactions(chain, walletAddress, category string, pageSize uint64) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.tatum.io/v4/data/transactions?chain=%s&addresses=%s&transactionSubtype=%s&pageSize=%d", chain, walletAddress, category, pageSize)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get user transactions")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetTransactionByHash(chain, hash string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.tatum.io/v4/data/transactions/hash?hash=%s&chain=%s", hash, chain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get user transactions")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
