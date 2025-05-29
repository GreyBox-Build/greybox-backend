package middlewares

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// LoadPublicKey loads the RSA public key from a file
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	readPath := strings.NewReader(path)
	keyData, err := io.ReadAll(readPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}

	return rsaPubKey, nil
}

// VerifySignature verifies the signature of the payload
func VerifySignature(publicKey *rsa.PublicKey, payload []byte, signature string) error {
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid base64 signature: %v", err)
	}

	hash := sha256.Sum256(payload)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		return errors.New("signature verification failed")
	}

	return nil
}

// SignatureMiddleware is the Gin middleware for verifying signatures
func SignatureMiddleware(publicKeyPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("x-signature")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing x-signature header"})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}

		c.Request.Body = io.NopCloser(nil)

		publicKey, err := LoadPublicKey(publicKeyPath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to load public key"})
			return
		}
		if err := VerifySignature(publicKey, body, signature); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
		
		c.Next()
	}
}
