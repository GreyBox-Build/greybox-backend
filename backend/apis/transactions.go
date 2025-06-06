package apis

import (
	"backend/serializers"
	"backend/state"
	"backend/utils/tokens"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TransactionRequest struct {
	Amount         string `json:"amount,omitempty"`
	Currency       string `json:"currency"`
	To             string `json:"to"`
	FeeCurrency    string `json:"feeCurrency"`
	FromPrivateKey string `json:"fromPrivateKey"` // Omit this field from JSON
}

type TransactionRequestV2 struct {
	FeeCurrency      string `json:"feeCurrency"`
	Chain            string `json:"chain"`
	CustodialAddress string `json:"custodialAddress"`
	Recipient        string `json:"recipient"`
	ContractType     uint32 `json:"contractType"`
	TokenAddress     string `json:"tokenAddress"`
	Amount           string `json:"amount"`
	FromPrivateKey   string `json:"fromPrivateKey"`
	Fee              struct {
		GasPrice string `json:"gasPrice"`
		GasLimit string `json:"gasLimit"`
	} `json:"fee"`
}

type ErrorResponse struct {
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       []Data `json:"data"`
}

type Data struct {
	Target      Target      `json:"target"`
	Value       int         `json:"value"`
	Property    string      `json:"property"`
	Constraints Constraints `json:"constraints"`
}

type Target struct {
	Property int `json:"property"`
}

type Constraints struct {
	Min string `json:"min"`
}

type PayoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		PayoutRequestID string `json:"payoutRequestId"`
		EscrowAddress   string `json:"escrowAddress"`
	} `json:"data"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type HurupayErrorResponse struct {
	Data    interface{}   `json:"data,omitempty"` // Accommodates an empty map or other data types
	Errors  []ErrorDetail `json:"errors,omitempty"`
	Message string        `json:"message"`
	Success bool          `json:"success"`
}

func GetUserTransactions(chain, walletAddress, category string, pageSize uint64) (map[string]interface{}, error) {
	url := ""
	switch category {
	case "":
		url = fmt.Sprintf("https://api.tatum.io/v4/data/transactions?chain=%s&addresses=%s&pageSize=%d", chain, walletAddress, pageSize)
	default:
		url = fmt.Sprintf("https://api.tatum.io/v4/data/transactions?chain=%s&addresses=%s&transactionSubtype=%s&pageSize=%d", chain, walletAddress, category, pageSize)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", state.AppConfig.TatumTestApiKey)

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

func GetTransactionByHash(chain, hash string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.tatum.io/v4/data/transactions/hash?hash=%s&chain=%s", hash, chain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", state.AppConfig.TatumTestApiKey)

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

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func PerformTransactionCelo(amount, accountAddress, privKey string, isNative bool) (string, int, error) {
	url := "https://api.tatum.io/v3/celo/transaction"
	client := &http.Client{}

	newData := TransactionRequest{
		Amount:         amount,
		Currency:       "CUSD",
		To:             accountAddress,
		FeeCurrency:    "CELO",
		FromPrivateKey: privKey,
	}
	if isNative {
		newData.Currency = "CELO"
	}

	// Convert the struct to JSON format
	jsonData, err := json.Marshal(&newData)
	if err != nil {
		return "", 500, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 500, err
	}

	req.Header.Add("x-api-key", state.AppConfig.TatumTestApiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", 500, err
	}
	defer resp.Body.Close()
	var errMsg string

	switch resp.StatusCode {
	case 403:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", 500, err
		}
		var result map[string]interface{}
		if err = json.Unmarshal(body, &result); err != nil {
			return "", 500, err
		}
		log.Println("result:", result)

		errMsg = "failed to perform transaction. most likely insufficient funds"

	case 400:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", 500, err
		}
		var result map[string]string
		if err = json.Unmarshal(body, &result); err != nil {
			return "", 500, err
		}
		errMsg = "failed to perform transaction. validation error"

	case 500:
		errMsg = "server error from third party application"
	case 401:
		errMsg = "subscription not active"
	default:
		// Handle any other status codes if needed
	}

	if errMsg != "" {
		return "", resp.StatusCode, errors.New(errMsg)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 500, err
	}

	var result map[string]string
	if err = json.Unmarshal(body, &result); err != nil {
		return "", 500, err
	}
	return result["txId"], resp.StatusCode, nil
}

func CalculateEstimatedFeeCelo(amount, to, from string) (map[string]interface{}, error) {
	apiUrl := "https://api.tatum.io/v3/celo/gas"
	client := &http.Client{}
	jsonData, err := json.Marshal(map[string]string{
		"amount": amount,
		"to":     to,
		"from":   from,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", state.AppConfig.TatumTestApiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to calculate estimated fee")
	}
	var result map[string]interface{}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func PerformTransactionXLM(data serializers.TransferXLM) (map[string]string, int, error) {
	apiUrl := "https://api.tatum.io/v3/xlm/transaction"
	client := &http.Client{}

	performRequest := func(data serializers.TransferXLM) (*http.Response, error) {
		requestData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestData))
		if err != nil {
			return nil, err
		}
		req.Header.Add("x-api-key", state.AppConfig.TatumTestApiKey)
		req.Header.Set("Content-type", "application/json")

		return client.Do(req)
	}

	resp, err := performRequest(data)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		data.Initialize = !data.Initialize
		resp, err = performRequest(data)
		if err != nil {
			return nil, 500, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 403 {
			errMsg := "failed to perform operation"
			failedResponse := map[string]interface{}{}
			if err := json.NewDecoder(resp.Body).Decode(&failedResponse); err != nil {
				return nil, 500, err
			}
			errorMessage, _ := failedResponse["message"].(string)
			mesage := map[string]string{
				"message": errorMessage,
			}
			return mesage, resp.StatusCode, fmt.Errorf(errMsg)
		}
	}

	errMsg := ""
	switch resp.StatusCode {
	case 200:
		respData := map[string]string{}
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return nil, 500, err
		}
		return respData, resp.StatusCode, nil
	case 400:
		errMsg = "validation error making request"
		return nil, 400, fmt.Errorf(errMsg)
	case 401:
		errMsg = "subscription not active anymore"
		return nil, 401, fmt.Errorf(errMsg)
	default:
		errMsg = "internal server error from Third Party application"
		return nil, 500, fmt.Errorf(errMsg)
	}
}

func GetUserTransactionXLM(address, pagination string) ([]serializers.TransactionXLM, error) {

	apiUrl := fmt.Sprintf("https://api.tatum.io/v3/xlm/account/tx/%s", address)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return []serializers.TransactionXLM{}, err
	}
	req.Header.Add("x-api-key", state.AppConfig.TatumTestApiKey)
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return []serializers.TransactionXLM{}, err
	}
	defer resp.Body.Close()

	log.Println("status code: ", resp.StatusCode)
	switch resp.StatusCode {
	case 200:
		respData := []serializers.TransactionXLM{}
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return []serializers.TransactionXLM{}, err
		}
		for _, data := range respData {
			tokens.DecodeXDR(data.EnvelopeXDR)
		}

		return respData, nil
	default:
		return []serializers.TransactionXLM{}, errors.New("failed to get user transactions")
	}
}

func DecodeTransactionDataXLM(data []serializers.TransactionXLM) (serializers.Result, error) {
	var t []serializers.Transaction
	for _, v := range data {
		txEnvelope, err := tokens.DecodeXDR(v.EnvelopeXDR)
		if err != nil {
			return serializers.Result{}, err
		}
		amount := tokens.GetAssociatedAmount(txEnvelope)
		tType := tokens.IsIncomingOrOutgoing(txEnvelope, v.ID)
		transType := tokens.GetTransactionAssetType(txEnvelope)
		tData := serializers.Transaction{
			Hash:               v.Hash,
			Chain:              "stellar",
			Address:            v.SourceAccount,
			BlockNumber:        v.Ledger,
			TransactionIndex:   0,
			TransactionType:    tType,
			Amount:             amount,
			Timestamp:          v.CreatedAt.Unix(),
			TokenAddress:       "",
			CounterAddress:     "",
			TransactionSubtype: transType,
		}
		t = append(t, tData)
	}
	result := serializers.Result{
		Result:   t,
		PrevPage: "0",
		NextPage: "0",
	}
	return result, nil
}

func GetTransactionByHashXLM(hash string) (serializers.TransactionXLM, error) {
	apiUrl := fmt.Sprintf("https://api.tatum.io/v3/xlm/transaction/%s", hash)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return serializers.TransactionXLM{}, err
	}
	req.Header.Add("x-api-key", state.AppConfig.TatumTestApiKey)
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return serializers.TransactionXLM{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return serializers.TransactionXLM{}, errors.New("failed to get user transactions")
	}

	respData := serializers.TransactionXLM{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return serializers.TransactionXLM{}, err
	}

	return respData, nil
}

type MobileMoneyResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    MobileData `json:"data"`
}

type MobileData struct {
	PartnerRequestID    string `json:"PartnerRequestID"`
	CollectionRequestID string `json:"CollectionRequestID"`
	ResponseCode        int    `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
}

type MobilePayoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		ResultCode        int    `json:"ResultCode"`
		PartnerRequestID  string `json:"PartnerRequestID"`
		ResultDescription string `json:"ResultDescription"`
	} `json:"data"`
}

func OnRampMobileMoney(data serializers.Payment) (MobileMoneyResponse, error) {
	apiUrl := "https://api.hurupay.com/v1/collections/mobile/initialize_transaction"
	client := &http.Client{}
	log.Println("onramp data: ", data)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return MobileMoneyResponse{}, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return MobileMoneyResponse{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", state.AppConfig.HurupayApiKey))
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return MobileMoneyResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errorResponse := HurupayErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return MobileMoneyResponse{}, err
		}
		log.Println("error response: ", errorResponse)
		return MobileMoneyResponse{}, errors.New(errorResponse.Message)
	}

	respData := MobileMoneyResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return MobileMoneyResponse{}, err
	}

	return respData, nil
}

func OffRampMobileMoney(data serializers.TransactionRequest) (PayoutResponse, error) {
	apiUrl := "https://api.hurupay.com/v1/payouts/mobile/initialize_transaction/request"
	client := &http.Client{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return PayoutResponse{}, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return PayoutResponse{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", state.AppConfig.HurupayApiKey))
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return PayoutResponse{}, err
	}
	defer resp.Body.Close()

	log.Println("status code: ", resp.StatusCode)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		errorResponse := HurupayErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("error response: ", errorResponse)
			return PayoutResponse{}, err
		}
		log.Println("error response: ", errorResponse)
		return PayoutResponse{}, errors.New(errorResponse.Message)
	}

	respData := PayoutResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return PayoutResponse{}, err
	}

	return respData, nil
}

func OffRampMobileFinalize(data serializers.TransactionDetails) (MobilePayoutResponse, error) {
	apiUrl := "https://api.hurupay.com/v1/payouts/mobile/initialize_transaction"
	client := &http.Client{}

	log.Println("data: ", data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return MobilePayoutResponse{}, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return MobilePayoutResponse{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", state.AppConfig.HurupayApiKey))
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return MobilePayoutResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		errorResponse := map[string]interface{}{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return MobilePayoutResponse{}, err
		}
		log.Println("error response finalize: ", errorResponse)
		return MobilePayoutResponse{}, errors.New("failed to perform transaction")
	}
	var output MobilePayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return MobilePayoutResponse{}, err
	}

	return output, nil
}
