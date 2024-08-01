package apis

import (
	"backend/models"
	"backend/serializers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func RegisterHmac() error {
	apiUrl := "https://api-eu1.tatum.io/v4/subscription"
	hmac := models.NewHmac()

	data := serializers.Hmac{HmacSecret: hmac.Secret}

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", apiUrl, bytes.NewBuffer(json))
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("content-type", "application/json")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to perform request which ended with status code: %d", resp.StatusCode)

	}
}

func CreateNotificationSubscription(address, chain string) error {
	apiUrl := "https://api-eu1.tatum.io/v4/subscription"
	webhookUrl := os.Getenv("WEBHOOK_URL")
	subType := os.Getenv("SUBSCRIPTION_TYPE")
	data := serializers.Subscription{Type: subType, Attr: serializers.SubscriptionAttr{Chain: chain, Url: webhookUrl, Address: address}}
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Add("x-api-key", os.Getenv("TATUM_API_KEY_TEST"))
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("status code: %d", resp.StatusCode)
	switch resp.StatusCode {
	case 201:
		return nil
	case 200:
		return nil
	default:
		return fmt.Errorf("failed to perform request which ended with status code: %d", resp.StatusCode)
	}
}
