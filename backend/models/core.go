package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"os"

	"gorm.io/gorm"
)

const IndexTo = 3

type ValidAddress struct {
	Index   int    `json:"index"`
	Address string `json:"address"`
}

type InvalidAddress struct {
	Address string `json:"address"`
	Reason  string `json:"reason"`
}

type AddressResponse struct {
	Valid   []ValidAddress   `json:"valid"`
	Invalid []InvalidAddress `json:"invalid"`
}

type MasterWallet struct {
	gorm.Model
	PublicAddress           string `json:"public_address"`
	PrivateKey              string `json:"private_key"`
	Mnemonic                string `json:"mnemonic"`
	XpublicAddress          string `json:"xpublic_address"`
	IsSetup                 bool   `json:"is_setup"`
	IndexFrom               uint64 `gorm:"default:0" json:"index_from"`
	IndexTo                 uint64 `gorm:"default:0" json:"index_to"`
	TotalAddressesGenerated uint64 `json:"total_addresses_generated"`
	WalletChain             string `gorm:"default:CELO" json:"wallet_chain"`
	TotalAddressActivated   uint64 `json:"total_address_activated"`
	SignatureId             string `json:"signature_id"`
}

type WalletAddress struct {
	gorm.Model
	PublicAddress  string       `json:"public_address"`
	IsActive       bool         `gorm:"default:false" json:"is_active"`
	WalletChain    string       `gorm:"default:CELO" json:"wallet_chain"`
	MasterWalletID uint         `json:"master_wallet_id"`
	MasterWallet   MasterWallet `gorm:"foreignKey:MasterWalletID;references:ID"`
	InUse          bool         `gorm:"default:false" json:"in_use"`
	WalletIndex    uint64       `json:"wallet_index"`
}

func (m *MasterWallet) UpdateMasterWallet() error {

	return db.Save(m).Error
}

func (w *WalletAddress) UpdateWalletAddress() error {

	return db.Save(w).Error
}

// an active address is an address that can receive funds
//
// this method activate the addresses that has been generated
func (m *MasterWallet) ActivatePumpAddreses() (map[string]interface{}, error) {
	apiUrl := "https://api.tatum.io/v3/gas-pump/activate"
	client := http.Client{}

	data := map[string]interface{}{
		"chain":          "CELO",
		"from":           m.IndexTo - IndexTo,
		"to":             m.IndexTo,
		"owner":          m.PublicAddress,
		"fromPrivateKey": m.PrivateKey,
	}
	// convert the data to json
	fmt.Println("req data:", data)
	// return nil, nil
	reqData, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := map[string]interface{}{}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
		fmt.Println(result)
		// return nil, fmt.Errorf("failed to activate pump addresses: %s, code", resp.Status)
		return nil, fmt.Errorf("failed to activate pump addresses: %s, code", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil

}

// this method is to return the lict of address that were successfully activated
func (m *MasterWallet) ResultAddressActivation(txId, chain string) (AddressResponse, error) {
	apiUrl := fmt.Sprintf("https://api.tatum.io/v3/gas-pump/address/%s/%s", chain, txId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return AddressResponse{}, err
	}
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	resp, err := client.Do(req)
	if err != nil {
		return AddressResponse{}, err
	}
	defer resp.Body.Close()
	result := AddressResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AddressResponse{}, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return AddressResponse{}, err
	}
	// update the master wallet
	m.TotalAddressActivated += uint64(len(result.Valid))
	if err := m.UpdateMasterWallet(); err != nil {
		return AddressResponse{}, err
	}

	return result, nil
}

func (m *MasterWallet) PrecalculatePumpAddresses() ([]string, int, error) {
	apiUrl := "https://api.tatum.io/v3/gas-pump"

	client := &http.Client{}

	var startIndex uint64
	if m.IndexTo == 0 {
		startIndex = m.IndexTo
	} else {
		startIndex = m.IndexTo + 1
	}

	data := map[string]interface{}{
		"chain": m.WalletChain,
		"owner": m.PublicAddress,
		"from":  startIndex,
		"to":    m.IndexTo + IndexTo,
	}
	m.IndexTo += IndexTo

	if err := m.UpdateMasterWallet(); err != nil {
		return nil, 500, err
	}
	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, 500, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, 500, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, err
	}
	errMsg := ""
	switch resp.StatusCode {
	case 200:
		var res []string

		if err := json.Unmarshal(body, &res); err != nil {
			return nil, 500, err
		}
		totalAddressCreated := len(res)
		m.TotalAddressesGenerated += uint64(totalAddressCreated)
		if err := m.UpdateMasterWallet(); err != nil {
			return nil, 500, err
		}
		return res, 200, nil
	case 400:
		errMsg = "Bad Request validation error"
	case 401:
		errMsg = "Unauthorized"
	case 500:
		errMsg = "Internal Server Error"
	default:
		errMsg = "Unknown Error"
	}
	return nil, resp.StatusCode, fmt.Errorf("%s", errMsg)

}

// fetch master wallet
func FetchMasterWallet(chain string) (MasterWallet, error) {
	var masterWallet MasterWallet
	if err := db.Where("wallet_chain = ?", chain).First(&masterWallet).Error; err != nil {
		return masterWallet, fmt.Errorf("failed to fetch master wallet: %w", err)
	}

	return masterWallet, nil
}

// BulkCreateWalletAddresses creates multiple WalletAddress records in the database
func BulkCreateWalletAddresses(publicAddresses []string, masterWalletID uint) error {
	var walletAddresses []WalletAddress

	for _, publicAddress := range publicAddresses {
		walletAddresses = append(walletAddresses, WalletAddress{
			PublicAddress:  publicAddress,
			MasterWalletID: masterWalletID,
			IsActive:       false, // Assuming new wallets are not active by default
			InUse:          false, // Assuming new wallets are not in use by default
		})
	}

	if len(walletAddresses) > 0 {
		if err := db.Create(&walletAddresses).Error; err != nil {
			return fmt.Errorf("failed to create wallet addresses: %w", err)
		}
	}

	return nil

}

// FindActiveAddress finds an active address from the database
func FindActiveAddress(chain string) (WalletAddress, error) {
	var address WalletAddress
	if err := db.Where("is_active = ? AND wallet_chain = ? AND in_use = ?", true, chain, false).First(&address).Error; err != nil {
		return address, fmt.Errorf("failed to find active address: %w", err)
	}
	// Mark the address as in use
	address.InUse = true
	if err := address.UpdateWalletAddress(); err != nil {
		return address, fmt.Errorf("failed to update address: %w", err)
	}

	return address, nil

}

// BulkUpdateWalletAddresses sets is_active to true for the given public addresses
func BulkUpdateWalletAddresses(publicAddresses []map[string]uint32) error {
	// Create a slice to hold the addresses
	var addresses []string
	// Create a map to hold the index values for each address
	addressIndexMap := make(map[string]uint32)

	// Extract the addresses and their indices
	for _, addrMap := range publicAddresses {
		for address, index := range addrMap {
			addresses = append(addresses, address)
			addressIndexMap[address] = index
		}
	}

	// Begin a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Prepare the updates in bulk
	for _, address := range addresses {
		if err := tx.Model(&WalletAddress{}).
			Where("public_address = ?", address).
			Updates(map[string]interface{}{"is_active": true, "wallet_index": addressIndexMap[address]}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update wallet address %s: %w", address, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetWalletAddress(publicAddress string) (WalletAddress, error) {
	var address WalletAddress
	if err := db.Where("public_address = ?", publicAddress).Preload("MasterWallet").First(&address).Error; err != nil {
		return address, fmt.Errorf("failed to get wallet address: %w", err)
	}

	return address, nil
}
