package utils

import (
	"backend/apis"
	"backend/models"
	"time"
)

// WebhookEvent represents the entire webhook payload
type WebhookEvent struct {
	Type           string    `json:"type"`
	OrganizationID string    `json:"organizationId"`
	Timestamp      int64     `json:"timestamp"`
	Data           EventData `json:"data"`
}

// EventData represents the transaction data
type EventData struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
	CreatedAt   time.Time   `json:"createdAt"`
	TxHash      *[]string   `json:"txHash"` // Nullable field
	FeeAmount   string      `json:"feeAmount"`
}

// Source represents the transaction source details
type Source struct {
	Amount       string `json:"amount"`
	FiatCurrency string `json:"fiatCurrency"`
}

// Destination represents the transaction destination details
type Destination struct {
	Asset     string `json:"asset"`
	AccountID string `json:"accountId"`
}

// BorderlessWebhookHandler handles all webhook related operations
func BorderlessWebhookHandler(body WebhookEvent) (bool, error) {

	requestId := body.Data.ID
	borderlessRequest, err := models.GetBorderlessRequestByTxId(requestId)
	if err != nil {
		return false, err
	}
	var trans models.Transaction
	trans.UserID = borderlessRequest.User.ID
	trans.User = borderlessRequest.User
	trans.Description = body.Type
	trans.RequestId = requestId
	trans.TransactionSubType = body.Data.Type
	trans.Amount = body.Data.Source.Amount
	trans.Asset = body.Data.Destination.Asset

	switch body.Type {
	case "Transaction_Created":
		borderlessRequest.Status = body.Data.Status
		if err := borderlessRequest.UpdateBorderlessRequest(); err != nil {
			return false, err
		}
		return true, nil
	case "Transaction_Updated":
		borderlessRequest.Status = body.Data.Status
		if body.Data.Status == "Completed" {
			// Process the transaction
			polygon := apis.NewTatumPolygon()
			masterWallet, err := models.FetchMasterWallet("MATIC")
			if err != nil {
				return false, err
			}
			hashResponse, err := polygon.PerformTransaction(
				borderlessRequest.User.AccountAddress, body.Data.Source.Amount,
				masterWallet.PrivateKey, apis.USDC_MATIC)
			if err != nil {
				return false, err
			}
			trans.CounterAddress = masterWallet.PublicAddress
			trans.Status = body.Data.Status
			trans.Hash = hashResponse.TxId
			trans.Chain = apis.MATIC.String()
			if err := trans.SaveTransaction(); err != nil {
				return false, err
			}
			nativeAmount, err := PerformDepositofNativeCalculation(trans.Amount, "USD", "MATIC")
			if err != nil {
				return false, err
			}
			go processMaticTransaction(*polygon, nativeAmount, &borderlessRequest.User, &masterWallet)

		} else {
			trans.Status = body.Data.Status
			if err := trans.SaveTransaction(); err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func processMaticTransaction(polygon apis.TatumPolygon, nativeAmount string, user *models.User, masterWallet *models.MasterWallet) {
	time.AfterFunc(1*time.Minute, func() {
		hashResponse, _ := polygon.PerformTransaction(
			user.AccountAddress, nativeAmount,
			masterWallet.PrivateKey, apis.MATIC)

		nativeTrans := createNativeTransaction(user, masterWallet, hashResponse.TxId, nativeAmount, "MATIC")
		_ = nativeTrans.SaveTransaction()
	})
}

func createNativeTransaction(user *models.User, masterWallet *models.MasterWallet, hash, amount, chain string) models.Transaction {
	return models.Transaction{
		Address:            user.AccountAddress,
		CounterAddress:     masterWallet.PublicAddress,
		Amount:             amount,
		UserID:             user.ID,
		User:               *user,
		Hash:               hash,
		Description:        "On-Ramp Deposit of Gas Fees",
		TransactionId:      hash,
		TransactionType:    "native",
		TransactionSubType: "Deposit",
		Chain:              chain,
		Asset:              chain,
		Status:             "Completed",
	}
}
