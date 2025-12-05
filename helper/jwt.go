package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SECRET KEY 
var SecretKey = []byte("SECRET_KEY")

// GENERATE TOKEN FULL
func GenerateFullToken(userID, roleID, studentID, lecturerID, roleName string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"role_id":     roleID,
		"role":        roleName,
		"student_id":  studentID,
		"lecturer_id": lecturerID,
		"exp":         time.Now().Add(24 * time.Hour).Unix(), // 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

//REFRESH & VALIDATE
func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
}
