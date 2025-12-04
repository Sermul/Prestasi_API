package helper

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("SECRET_KEY")

func GenerateToken(userID string, roleID string, studentID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":    userID,
        "role_id":    roleID,
        "student_id": studentID, // boleh "" kalau bukan mahasiswa
        "exp":        time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(SecretKey)
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
    return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        return SecretKey, nil
    })
}
