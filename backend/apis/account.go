package apis

import (
	"backend/models"
	"backend/serializers"
	"backend/state"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

func SetupCeloAccount(user *models.User) error {
	mnemonic, xpub, err := GenerateCelloWallet()
	if err != nil {
		return fmt.Errorf("generate cello wallet: %w", err)
	}

	address, err := GenerateCelloAddress(xpub)
	if err != nil {
		return fmt.Errorf("generate cello address: %w", err)
	}

	privData := serializers.PrivGeneration{Index: 1, Mnemonic: mnemonic}
	privateKey, err := GeneratePrivateKey(privData)
	if err != nil {
		return fmt.Errorf("generate cello private key: %w", err)
	}

	encryptedKey, err := encrypt(privateKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	encryptedMnemonic, err := encrypt(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}

	encrytedXpub, err := encrypt(xpub)
	if err != nil {
		return fmt.Errorf("failed to encrypt xpub: %w", err)
	}

	user.Mnemonic = encryptedMnemonic
	user.Xpub = encrytedXpub
	user.AccountAddress = address
	user.PrivateKey = encryptedKey
	return nil
}

func SetupStellarAccount(user *models.User) error {
	data, _, err := GenerateXlmAccount()
	if err != nil {
		return fmt.Errorf("generate stellar account: %w", err)
	}
	user.AccountAddress = data["address"]

	encryptedSecret, err := encrypt(data["secret"])
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	user.PrivateKey = encryptedSecret
	return nil
}

func SetupPolygonAccount(user *models.User) error {
	polygon := NewTatumPolygon()

	wallet, err := polygon.CreateWallet()
	if err != nil {
		return fmt.Errorf("create polygon wallet: %w", err)
	}

	privResp, err := polygon.GeneratePrivateKey(wallet.Mnemonic, 0)
	if err != nil {
		return fmt.Errorf("generate polygon private key: %w", err)
	}

	addrResp, err := polygon.GenerateAddress(wallet.Xpub, 0)
	if err != nil {
		return fmt.Errorf("generate polygon address: %w", err)
	}

	encryptedKey, err := encrypt(privResp.Key)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	encryptedMnemonic, err := encrypt(wallet.Mnemonic)
	if err != nil {
		return fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}

	encryptedXpub, err := encrypt(wallet.Xpub)
	if err != nil {
		return fmt.Errorf("failed to encrypt xpub: %w", err)
	}

	user.Mnemonic = encryptedMnemonic
	user.Xpub = encryptedXpub
	user.PrivateKey = encryptedKey
	user.AccountAddress = addrResp.Address
	return nil
}

func encrypt(plainText string) (string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(state.AppConfig.EncryptionKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(cipherTextB64 string) (string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(state.AppConfig.EncryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
