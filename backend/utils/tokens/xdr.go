package tokens

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/stellar/go/xdr"
)

func DecodeXDR(xdrString string) (*xdr.TransactionEnvelope, error) {
	// Decode the base64 XDR string
	raw, err := base64.StdEncoding.DecodeString(xdrString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %v", err)
	}

	// Decode the raw XDR data into a TransactionEnvelope
	newRaw := bytes.NewReader(raw)
	// var txEnvelope xdr.TransactionEnvelope
	var txEnvelope xdr.TransactionEnvelope
	if _, err := xdr.Unmarshal(newRaw, &txEnvelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XDR: %v", err)
	}
	return &txEnvelope, nil
}

func GetAssociatedAmount(txEnvelope *xdr.TransactionEnvelope) string {
	tx := txEnvelope.V1.Tx
	for _, op := range tx.Operations {
		switch op.Body.Type {
		case xdr.OperationTypePayment:
			paymentOp := op.Body.MustPaymentOp()
			return fmt.Sprintf("%d", paymentOp.Amount)
		case xdr.OperationTypePathPaymentStrictReceive:
			pathPaymentOp := op.Body.MustPathPaymentStrictReceiveOp()
			return fmt.Sprintf("%d", pathPaymentOp.DestAmount)
		case xdr.OperationTypeCreateAccount:
			createAccountOp := op.Body.MustCreateAccountOp()
			return fmt.Sprintf("%d", createAccountOp.StartingBalance)
		// Add other operation types as needed
		default:
			return "0"
		}

	}
	return "0"
}

// IsIncomingOrOutgoing determines if the transaction is incoming or outgoing for the given account
func IsIncomingOrOutgoing(txEnvelope *xdr.TransactionEnvelope, accountID string) string {
	tx := txEnvelope.V1.Tx
	for _, op := range tx.Operations {
		switch op.Body.Type {
		case xdr.OperationTypePayment:
			paymentOp := op.Body.MustPaymentOp()
			if string(paymentOp.Destination.Address()) == accountID {
				return "incoming"
			}
			if op.SourceAccount != nil && string(op.SourceAccount.Address()) == accountID {
				return "outgoing"
			}
		case xdr.OperationTypePathPaymentStrictReceive:
			pathPaymentOp := op.Body.MustPathPaymentStrictReceiveOp()
			if string(pathPaymentOp.Destination.Address()) == accountID {
				return "incoming"
			}
			if op.SourceAccount != nil && string(op.SourceAccount.Address()) == accountID {
				return "outgoing"
			}
		case xdr.OperationTypeCreateAccount:
			createAccountOp := op.Body.MustCreateAccountOp()
			if string(createAccountOp.Destination.Address()) == accountID {
				return "incoming"
			}
			if op.SourceAccount != nil && string(op.SourceAccount.Address()) == accountID {
				return "outgoing"
			}

		}
	}
	return "unknown"
}

// GetTransactionAssetType determines the asset type involved in the transaction
func GetTransactionAssetType(txEnvelope *xdr.TransactionEnvelope) string {
	tx := txEnvelope.V1.Tx
	for _, op := range tx.Operations {
		switch op.Body.Type {
		case xdr.OperationTypePayment:
			paymentOp := op.Body.MustPaymentOp()
			return getAssetType(paymentOp.Asset)
		case xdr.OperationTypePathPaymentStrictReceive:
			pathPaymentOp := op.Body.MustPathPaymentStrictReceiveOp()
			return getAssetType(pathPaymentOp.SendAsset)
		case xdr.OperationTypePathPaymentStrictSend:
			pathPaymentOp := op.Body.MustPathPaymentStrictSendOp()
			return getAssetType(pathPaymentOp.SendAsset)
		case xdr.OperationTypeCreateAccount:
			return "native" // CreateAccount operations involve native Lumens
		default:
			return "unknown"
		}
	}
	return "unknown"
}

// getAssetType determines the type of a given asset
func getAssetType(asset xdr.Asset) string {
	switch asset.Type {
	case xdr.AssetTypeAssetTypeNative:
		return "native"
	case xdr.AssetTypeAssetTypeCreditAlphanum4, xdr.AssetTypeAssetTypeCreditAlphanum12:
		return "fungible"
	default:
		return "unknown"
	}
}
