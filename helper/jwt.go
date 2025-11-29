package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("SECRET_KEY")

// Generate JWT Token
func GenerateToken(userID string, roleID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role_id": roleID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SecretKey)
}

// Parse & Validate Token → dipakai middleware
func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
}
