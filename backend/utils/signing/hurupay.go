package signing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func VerifyWebhookSignature(body string, signature string, publicKeyStr string) (bool, error) {
	// Decode the PEM-encoded public key string
	block, _ := pem.Decode([]byte(publicKeyStr))

	if block == nil {
		return false, errors.New("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not an RSA public key")
	}

	// Hash the body content using SHA-256
	hashedData := sha256.Sum256([]byte(body))

	// Decode the base64-encoded signature
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	// Verify the signature using the RSA public key
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashedData[:], decodedSignature)
	if err != nil {
		return false, err
	}

	return true, nil
}
