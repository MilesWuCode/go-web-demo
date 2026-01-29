package server

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponses)
}

// GetUserByID 處理 GET /api/users/{id} 請求
func (app *Application) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// 從 URL 路徑中獲取 id
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// 使用 GORM Gen 的型別安全查詢
	var user models.User
	if err := app.DB.First(&user, id).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	// 建立並回傳安全的回應
	userResponse := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}
