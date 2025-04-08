package models

import (
	"fmt"

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
	PrivateKey              string `json:"-"`
	Mnemonic                string `json:"-"`
	XpublicAddress          string `json:"-"`
	IsSetup                 bool   `json:"is_setup"`
	CurrentIndex            uint64 `json:"currrent_index"`
	TotalAddressesGenerated uint64 `json:"-"`
	WalletChain             string `gorm:"default:CELO" json:"wallet_chain"`
	TotalAddressActivated   uint64 `json:"-"`
	SignatureId             string `json:"-"`
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

func (m *MasterWallet) CreateMasterWallet() error {
	return db.Create(m).Error
}

// fetch master wallet
func FetchMasterWallet(chain string) (MasterWallet, error) {
	var masterWallet MasterWallet
	if err := db.Where("wallet_chain = ?", chain).Last(&masterWallet).Error; err != nil {
		return masterWallet, fmt.Errorf("failed to fetch master wallet: %w", err)
	}

	return masterWallet, nil
}

// fetch all master wallets
func FetchMasterWallets() ([]MasterWallet, error) {
	var masterWallets []MasterWallet
	if err := db.Find(&masterWallets).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch master wallets: %w", err)
	}
	return masterWallets, nil

}
