package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func WebhookSignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
			fmt.Println("Unable to read request body")
			return
		}

		// Reassign the body to allow further reads down the middleware chain
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Get the signature from the header
		signature := c.GetHeader("x-webhook-signature")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing signature"})
			fmt.Println("Missing signature")
			return
		}
		dir, _ := os.Getwd()
		// Determine the PEM file path based on the webhook path
		var pemFilePath string
		switch c.Request.URL.Path {
		case "/api/v1/notification/on-ramp":
			pemFilePath = dir + "/onramp-public.pem"
		case "/api/v1/notification/off-ramp":
			pemFilePath = dir + "/offramp-public.pem"
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook path"})
			fmt.Println("Invalid webhook path")
			return
		}

		// Read the PEM file using os.ReadFile
		_, err = os.ReadFile(pemFilePath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to read public key file"})
			fmt.Println("Unable to read public key file", err)
			return
		}

		fmt.Println("signature", signature)

		// Verify the webhook signature
		//flag, err := signing.VerifyWebhookSignature(bodyString, signature, publicKey)
		//if err != nil || !flag {
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		//	fmt.Println("Invalid signature", err)
		//		return
		//	}

		// If verification is successful, proceed to the next handler
		c.Next()
	}
}
