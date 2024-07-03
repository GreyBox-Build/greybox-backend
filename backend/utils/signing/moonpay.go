package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
)

func GenerateSignedURL(originalURL, secretKey string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(parsedURL.RawQuery))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	signedURL := originalURL + "&signature=" + url.QueryEscape(signature)
	return signedURL, nil
}
