package api

import (
	"errors"
	"net/http"
	"time"
	"web-demo/model"
	"web-demo/server"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
)

// RegisterRequest 定義了使用者註冊時預期的請求資料格式。
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest 定義了使用者登入時預期的請求資料格式。
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// AuthHandler 結構體，用於處理認證相關的 API 請求
type AuthHandler struct {
	App *server.Application
}

// NewAuthHandler 建立並回傳一個新的 AuthHandler 實例
func NewAuthHandler(app *server.Application) *AuthHandler {
	return &AuthHandler{App: app}
}

// Register 處理使用者註冊請求
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterRequest

	// 解析請求主體
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// 驗證請求資料
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		h.App.ErrorJSON(w, validationErrors, http.StatusBadRequest)
		return
	}

	// 檢查電子郵件是否已存在
	var existingUser model.User
	result := h.App.DB.Where("email = ?", input.Email).First(&existingUser)
	if result.Error == nil { // Found existing user
		h.App.ErrorJSON(w, errors.New("此電子郵件已被註冊"), http.StatusConflict) // 409 Conflict
		return
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // Other database error
		h.App.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	// 雜湊密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		h.App.ErrorJSON(w, errors.New("無法雜湊密碼"), http.StatusInternalServerError)
		return
	}

	// 建立新使用者
	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// 儲存使用者到資料庫
	result = h.App.DB.Create(&user)
	if result.Error != nil {
		h.App.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	// 回傳成功回應 (不包含密碼)
	h.App.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "使用者註冊成功",
		"user": UserResponse{ // 使用現有的 UserResponse 結構
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// Login 處理使用者登入請求
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest

	// 解析請求主體
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// 驗證請求資料
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		h.App.ErrorJSON(w, validationErrors, http.StatusBadRequest)
		return
	}

	// 檢查使用者是否存在
	var user model.User
	result := h.App.DB.Where("email = ?", input.Email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		h.App.ErrorJSON(w, errors.New("無效的憑證"), http.StatusUnauthorized)
		return
	}
	if result.Error != nil {
		h.App.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	// 比對密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		h.App.ErrorJSON(w, errors.New("無效的憑證"), http.StatusUnauthorized)
		return
	}

	// 產生 JWT token
	token, err := h.generateToken(&user)
	if err != nil {
		h.App.ErrorJSON(w, errors.New("無法產生認證 token"), http.StatusInternalServerError)
		return
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "登入成功",
		"token":   token,
		"user": UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// generateToken 根據使用者資料產生 JWT token
func (h *AuthHandler) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * time.Duration(h.App.Config.JWTExpiresIn)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.App.Config.JWTSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// LogoutRequest 定義了使用者登出時預期的請求資料格式。
type LogoutRequest struct {
	Token string `json:"token" validate:"required"`
}

// Logout 處理使用者登出請求
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var input LogoutRequest

	// 解析請求主體
	err := h.App.ReadJSON(w, r, &input)
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// 驗證 token 是否存在 (不需要進一步驗證有效性，因為假設客戶端處理銷毀)
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		h.App.ErrorJSON(w, validationErrors, http.StatusBadRequest)
		return
	}

	// 這裡僅回傳成功訊息，實際的 token 清除由客戶端處理。
	// 若要實作伺服器端 token 無效化 (例如 JWT 黑名單)，則需更複雜的機制。
	h.App.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "登出成功",
	})
}

