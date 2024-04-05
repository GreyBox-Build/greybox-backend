package tokens

import (
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

var (
	apiSecret = []byte(os.Getenv("API_SECRET"))
)

func GenerateToken(userID uint) (string, error) {
	claims := CustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Example: 1 day expiration
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(apiSecret)
}

func IsTokenValid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return apiSecret, nil
	})

	return err == nil && token.Valid
}

func ExtractToken(c *gin.Context) string {
	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		// Token not found in the header, check query parameter or cookie if needed
		authHeader = c.Query("token")
	}

	// Token could not be found in the header or query parameter
	if authHeader == "" {
		return ""
	}

	// Check if the header has the correct format ("Bearer <token>")
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return ""
	}

	return splitToken[1]
}

func ExtractUserID(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return apiSecret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse JWT token: %v", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid or expired token")
	}

	return claims.ID, nil
}
