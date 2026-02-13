package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apiHandler "web-demo/handler/api"
	"web-demo/model"
	"web-demo/server"
	utilAuth "web-demo/util/auth"

	"google.golang.org/api/idtoken"
)

// GoogleVerifyRequest 定義了從前端接收的 Google 身份憑證資料格式
type GoogleVerifyRequest struct {
	IDToken string `json:"id_token"`
}

// GoogleVerifyHandler 結構體，用於處理 Google 身份憑證驗證相關的 API 請求
type GoogleVerifyHandler struct {
	App *server.Application
}

// NewGoogleVerifyHandler 建立並回傳一個新的 GoogleVerifyHandler 實例
func NewGoogleVerifyHandler(app *server.Application) *GoogleVerifyHandler {
	return &GoogleVerifyHandler{App: app}
}

// VerifyGoogleIDToken 處理 Google ID Token 驗證請求
func (h *GoogleVerifyHandler) VerifyGoogleIDToken(w http.ResponseWriter, r *http.Request) {
	var input GoogleVerifyRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無效的請求資料: %v", err), http.StatusBadRequest)
		return
	}

	// 驗證 Google ID Token
	payload, err := idtoken.Validate(context.Background(), input.IDToken, h.App.Config.GoogleClientID)
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("Google ID Token 驗證失敗: %v", err), http.StatusUnauthorized)
		return
	}

	// 從 payload 中提取用戶資訊
	googleID := payload.Subject // Google ID
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	// 檢查 SocialProvider 資料表
	var socialProvider model.SocialProvider
	err = h.App.DB.Where("provider = ? AND provider_id = ?", "google", googleID).First(&socialProvider).Error

	var user model.User
	if err != nil { // SocialProvider 不存在
		// 檢查是否已存在相同 email 的用戶
		err = h.App.DB.Where("email = ?", email).First(&user).Error
		if err != nil { // 用戶不存在，創建新用戶
			user = model.User{
				Name:  name,
				Email: email,
			}
			if err := h.App.DB.Create(&user).Error; err != nil {
				h.App.ErrorJSON(w, fmt.Errorf("無法創建用戶: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// 創建 SocialProvider 記錄
		socialProvider = model.SocialProvider{
			UserID:     user.ID,
			Provider:   "google",
			ProviderID: googleID,
			Email:      email,
			Name:       name,
		}
		if err := h.App.DB.Create(&socialProvider).Error; err != nil {
			h.App.ErrorJSON(w, fmt.Errorf("無法創建 SocialProvider 記錄: %v", err), http.StatusInternalServerError)
			return
		}
	} else { // SocialProvider 存在，獲取關聯的用戶
		err = h.App.DB.First(&user, socialProvider.UserID).Error
		if err != nil {
			h.App.ErrorJSON(w, fmt.Errorf("無法找到關聯用戶: %v", err), http.StatusInternalServerError)
			return
		}
		// 如果用戶名稱或 email 不一致，可以選擇更新
		if user.Name == "" || user.Name != name {
			user.Name = name
		}
		if user.Email != email {
			user.Email = email
		}
		if err := h.App.DB.Save(&user).Error; err != nil {
			h.App.ErrorJSON(w, fmt.Errorf("無法更新用戶資訊: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// 產生 JWT token
	appToken, err := utilAuth.GenerateToken(&user, time.Minute*time.Duration(h.App.Config.JWTExpiresIn), h.App.Config.JWTSecret)
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法產生應用程式 token: %v", err), http.StatusInternalServerError)
		return
	}

	// 產生 Refresh Token (如果應用需要)
	refreshTokenString, err := utilAuth.CreateAndStoreRefreshToken(h.App.DB, user.ID, time.Minute*time.Duration(h.App.Config.JWTRefreshExpiresIn))
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法產生換新權杖: %v", err), http.StatusInternalServerError)
		return
	}

	// 回傳 token 和使用者資訊
	h.App.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": apiHandler.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		"access_token":             appToken,
		"refresh_token":            refreshTokenString,
		"access_token_expires_in":  h.App.Config.JWTExpiresIn * 60, // Convert minutes to seconds
		"refresh_token_expires_in": h.App.Config.JWTRefreshExpiresIn * 60,
	})
}
