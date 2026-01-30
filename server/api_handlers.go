package server

import (
	"errors"
	"net/http"
	"strconv"
	"web-demo/models"

	"gorm.io/gorm"
)

// --- API Handlers ---

func (app *Application) EchoHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"data": "hello world"}
	app.writeJSON(w, http.StatusOK, data)
}

// UserResponse 用於定義回傳給前端的使用者資料結構，以隱藏密碼等敏感資訊
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetAllUsers 處理 GET /api/users 請求
func (app *Application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	// 使用 WithContext 傳遞請求的 context
	if err := app.DB.WithContext(r.Context()).Model(&models.User{}).Find(&users).Error; err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	app.writeJSON(w, http.StatusOK, userResponses)
}

// GetUserByID 處理 GET /api/users/{id} 請求
func (app *Application) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64) // 使用 ParseUint 確保 ID 為正整數
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID"), http.StatusBadRequest)
		return
	}

	var user models.User
	// 使用 WithContext 傳遞請求的 context
	err = app.DB.WithContext(r.Context()).First(&user, id).Error
	if err != nil {
		// 判斷是否為「找不到紀錄」的特定錯誤
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.errorJSON(w, errors.New("user not found"), http.StatusNotFound)
		} else {
			app.errorJSON(w, err, http.StatusInternalServerError)
		}
		return
	}

	userResponse := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	app.writeJSON(w, http.StatusOK, userResponse)
}
