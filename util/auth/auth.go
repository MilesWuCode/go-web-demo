package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"
	"web-demo/model"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken 根據使用者資料產生 JWT token
func GenerateToken(user *model.User, jwtExpiresIn time.Duration, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(jwtExpiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// GenerateRefreshToken 生成一個安全的隨機字串作為換新權杖
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
