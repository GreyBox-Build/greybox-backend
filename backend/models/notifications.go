package models

import (
	"fmt"
	"time"

	"backend/serializers"

	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"os"

	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hmac struct {
	*gorm.Model
	Secret string `gorm:"type:varchar(64)" json:"hmac_secret"`
}

type Notification struct {
	*gorm.Model
	UserID uint      `json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
	Type   string    `json:"type"`
	Body   string    `json:"body"`
	Read   bool      `json:"read"`
	Link   string    `json:"link"`
	ReadAt time.Time `json:"read_at"`
}

func (n *Notification) MarkAsRead() {
	n.Read = true
	n.ReadAt = time.Now()
	db.Save(n)

}

func (n *Notification) CreateNotification() error {
	return db.Create(&n).Error
}

// GenerateHmacSecret generates a new HMAC secret
func GenerateHmacSecret() (string, error) {
	// Generate a random key
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate random key: %v", err)
	}

	// Generate a new HMAC using SHA256
	h := hmac.New(sha256.New, key)
	message := []byte(uuid.New().String()) // Using a UUID as the message for uniqueness
	h.Write(message)

	// Encode the HMAC in hexadecimal
	secret := hex.EncodeToString(h.Sum(nil))
	return secret, nil
}

// NewHmac creates a new Hmac instance with a generated secret
func NewHmac() *Hmac {
	secret, err := GenerateHmacSecret()
	if err != nil {
		fmt.Printf("Error generating HMAC secret: %v\n", err)
		return nil
	}

	h := &Hmac{
		Secret: secret,
	}
	fmt.Println("HMAC secret: ", h.Secret)
	return h
}

func VerifyWebhookAuthenticity(secret string, payload serializers.Webhook) bool {
	hmacSecret := os.Getenv("HMAC_SECRET")

	Json, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload: ", err)
		return false
	}
	h := hmac.New(sha512.New, []byte(hmacSecret))
	h.Write(Json)
	base64Hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	checkValues := secret == base64Hash

	return checkValues
}

func GetUserNotifications(user *User, isRead bool) []Notification {
	var notifications []Notification
	db.Where("user_id = ? AND read = ?", user.ID, isRead).Find(&notifications)
	return notifications
}
