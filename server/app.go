package server

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"time"
	"web-demo/config"
	"web-demo/model"

	"gorm.io/gorm"
)

// Application 結構用來存放共享的依賴
type Application struct {
	Config        *config.AppConfig
	DB            *gorm.DB
	TemplateCache map[string]*template.Template
}

// GetJWTSecret 取得 JWT Secret
func (app *Application) GetJWTSecret() string {
	return app.Config.JWTSecret
}

// GenerateRefreshToken 產生並儲存一個新的換新權杖
func (app *Application) GenerateRefreshToken(userID uint, expiry time.Time) (string, error) {
	b := make([]byte, 32) // 建立一個 32 位元組的隨機字節切片
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	refreshTokenString := base64.URLEncoding.EncodeToString(b) // 將隨機字節編碼為 URL 安全的 base64 字串

	refreshToken := model.RefreshToken{
		UserID:    userID,
		Token:     refreshTokenString,
		ExpiresAt: expiry,
	}

	if err := app.DB.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return refreshTokenString, nil
}
