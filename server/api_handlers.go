package server

import (
	"encoding/json"
	"net/http"
	"web-demo/models"
)

// --- API Handlers ---

func (app *Application) EchoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{"data": "hello world"}
	json.NewEncoder(w).Encode(data)
}

// UserResponse 用於定義回傳給前端的使用者資料結構，以隱藏密碼等敏感資訊
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetAllUsers 處理 GET /api/users 請求
func (app *Application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// 使用 GORM Gen 的型別安全查詢
	var users []models.User
	if err := app.DB.Model(&models.User{}).Find(&users).Error; err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
		return
	}

	// 建立一個新的 slice 來存放處理過的 user 資料
	// 這是為了避免回傳密碼等敏感資訊
	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	// 設定回應標頭並回傳 JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponses)
}
