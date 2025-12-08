package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("YOUR_SUPER_SECRET_KEY") // ⚡ Move to env variable in prod!

// GenerateToken generates a JWT token for a given userID
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),                     // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}
