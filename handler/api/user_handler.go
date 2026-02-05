package api

import (
	"errors"
	"net/http"
	"strconv"
	"web-demo/generated/query"
	"web-demo/util"

	"gorm.io/gorm"
)

// PaginatedUserResponse 是包含分頁資訊的使用者列表回應結構
type PaginatedUserResponse struct {
	Data       []UserResponse  `json:"data"`
	Pagination util.Pagination `json:"pagination"`
}

// GetAllUsers 處理 GET /api/users 請求，並支援分頁
func (h *APIHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// --- 讀取分頁參數 ---
	page := h.App.ReadIntQuery(r, "page", 1)
	pageSize := h.App.ReadIntQuery(r, "pageSize", 10)

	// --- 查詢資料庫 ---
	q := query.Use(h.App.DB)
	u := q.User

	// 獲取總記錄數
	totalRecords, err := u.WithContext(r.Context()).Count()
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// 使用 Limit 和 Offset 實現分頁查詢
	users, err := u.WithContext(r.Context()).Limit(pageSize).Offset((page - 1) * pageSize).Find()
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// --- 轉換為回應結構 ---
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	// --- 組合分頁回應 ---
	paginatedResponse := PaginatedUserResponse{
		Data:       userResponses,
		Pagination: util.NewPagination(page, pageSize, totalRecords),
	}

	// 回傳 JSON 資料
	h.App.WriteJSON(w, http.StatusOK, paginatedResponse)
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
