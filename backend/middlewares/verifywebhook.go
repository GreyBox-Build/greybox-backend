package middlewares

import (
	"backend/utils/signing"
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func WebhookSignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		bodyString := string(bodyBytes)

		// Get the signature from the header
		signature := c.GetHeader("x-webhook-signature")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing signature"})
			return
		}

		var secret string
		switch c.Request.URL.Path {
		case "/api/v1/notification/on-ramp":
			secret = os.Getenv("WEBHOOK_ONRAMP_PUBLIC_KEY")
		case "//api/v1/notification/off-ramp":
			secret = os.Getenv("WEBHOOK_OFFRAMP_PUBLIC_KEY")
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook path"})
			return
		}

		// Verify the webhook signature
		flag, err := signing.VerifyWebhookSignature(bodyString, signature, secret)
		if err != nil || !flag {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}

		// If verification is successful, proceed to the next handler
		c.Next()
	}
}
