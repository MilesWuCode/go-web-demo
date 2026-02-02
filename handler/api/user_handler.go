package api

import (
	"errors"
	"net/http"
	"strconv"
	"web-demo/generated/query"

	"gorm.io/gorm"
)

// UserResponse 用於定義回傳給前端的使用者資料結構，以隱藏密碼等敏感資訊
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetAllUsers 處理 GET /api/users 請求
func (h *APIHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	q := query.Use(h.App.DB)
	users, err := q.User.WithContext(r.Context()).Find()
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusInternalServerError)
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

	h.App.WriteJSON(w, http.StatusOK, userResponses)
}

// GetUserByID 處理 GET /api/users/{id} 請求
func (h *APIHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64) // 使用 ParseUint 確保 ID 為正整數
	if err != nil {
		h.App.ErrorJSON(w, errors.New("invalid user ID"), http.StatusBadRequest)
		return
	}

	q := query.Use(h.App.DB)
	u := q.User
	user, err := u.WithContext(r.Context()).Where(u.ID.Eq(uint(id))).First()
	if err != nil {
		// 判斷是否為「找不到紀錄」的特定錯誤
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.App.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		} else {
			h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		}
		return
	}

	userResponse := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	h.App.WriteJSON(w, http.StatusOK, userResponse)
}

// GetUserByName 處理 GET /api/username/{name} 請求
func (h *APIHandler) GetUserByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		h.App.ErrorJSON(w, errors.New("invalid user name"), http.StatusBadRequest)
		return
	}

	q := query.Use(h.App.DB)
	user, err := q.User.WithContext(r.Context()).FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.App.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		} else {
			h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		}
		return
	}

	userResponse := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	h.App.WriteJSON(w, http.StatusOK, userResponse)
}
