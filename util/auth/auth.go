package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
	"web-demo/model"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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

// CreateAndStoreRefreshToken 生成一個新的 refresh token，存入資料庫並返回
func CreateAndStoreRefreshToken(db *gorm.DB, userID uint, expiresIn time.Duration) (string, error) {
	// 1. 生成隨機 token 字串
	token, err := GenerateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("無法生成 refresh token 字串: %w", err)
	}

	// 2. 計算過期時間
	expiresAt := time.Now().Add(expiresIn)

	// 3. 建立紀錄
	refreshToken := model.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	// 4. 存入資料庫
	if result := db.Create(&refreshToken); result.Error != nil {
		return "", fmt.Errorf("無法儲存 refresh token: %w", result.Error)
	}

	return token, nil
}
