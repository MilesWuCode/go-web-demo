package auth

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	apiHandler "web-demo/handler/api"
	"web-demo/middleware"
	"web-demo/model"
	"web-demo/server"
	"web-demo/util/auth"

	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh_tw"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterRequest 定義了使用者註冊時預期的請求資料格式。
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50" label:"姓名"`
	Email    string `json:"email" validate:"required,email" label:"信箱"`
	Password string `json:"password" validate:"required,min=6" label:"密碼"`
}

// LoginRequest 定義了使用者登入時預期的請求資料格式。
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

var customMessages = map[string]string{
	"required": "此欄位為必填",
	"email":    "請輸入有效的電子郵件格式",
	"min":      "長度不足",
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
	err := h.App.ParseRequestForm(w, r)
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Manually populate input from form values
	input.Name = r.PostForm.Get("name")
	input.Email = r.PostForm.Get("email")
	input.Password = r.PostForm.Get("password")

	// 驗證請求資料
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 優先使用 label tag
		if name := fld.Tag.Get("label"); name != "" {
			return name
		}

		return fld.Name
	})
	zhTW := zh_Hant_TW.New()
	uni := ut.New(zhTW, zhTW)
	trans, _ := uni.GetTranslator("zh_Hant_TW")
	zh_translations.RegisterDefaultTranslations(validate, trans)
	err = validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			// validationErrors[err.Field()] = err.Tag()
			// validationErrors[err.Field()] = customMessages[err.Tag()]
			validationErrors[err.Field()] = err.Translate(trans)
		}

		h.App.WriteJSON(w, http.StatusUnprocessableEntity, validationErrors)
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
		"user": apiHandler.UserResponse{ // 使用現有的 UserResponse 結構
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
	err := h.App.ParseRequestForm(w, r)
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Manually populate input from form values
	input.Email = r.PostForm.Get("email")
	input.Password = r.PostForm.Get("password")

	// 驗證請求資料
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}

		h.App.WriteJSON(w, http.StatusUnprocessableEntity, validationErrors)
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
	token, err := auth.GenerateToken(&user, time.Minute*time.Duration(h.App.Config.JWTExpiresIn), h.App.Config.JWTSecret)
	if err != nil {
		h.App.ErrorJSON(w, errors.New("無法產生認證 token"), http.StatusInternalServerError)
		return
	}

	// 產生並儲存 Refresh Token
	refreshToken, err := auth.CreateAndStoreRefreshToken(h.App.DB, user.ID, time.Minute*time.Duration(h.App.Config.JWTRefreshExpiresIn))
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法建立換新 token: %w", err), http.StatusInternalServerError)
		return
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "登入成功",
		"user": apiHandler.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		"access_token":             token,
		"refresh_token":            refreshToken,
		"access_token_expires_in":  h.App.Config.JWTExpiresIn * 60, // Convert minutes to seconds
		"refresh_token_expires_in": h.App.Config.JWTRefreshExpiresIn * 60,
	})
}

// LogoutRequest 定義了使用者登出時預期的請求資料格式。
type LogoutRequest struct {
	Token string `json:"token" validate:"required"`
}

// Logout 處理使用者登出請求
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var input LogoutRequest

	// 解析請求主體
	err := h.App.ParseRequestForm(w, r)
	if err != nil {
		h.App.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Manually populate input from form values
	input.Token = r.PostForm.Get("token")

	// 驗證 token 是否存在 (不需要進一步驗證有效性，因為假設客戶端處理銷毀)
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		h.App.WriteJSON(w, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	// 清除授權：刪除指定的 Refresh Token
	// 這裡假設傳入的 token 是要被撤銷的 refresh token
	if err := h.App.DB.Where("token = ?", input.Token).Delete(&model.RefreshToken{}).Error; err != nil {
		h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "登出成功",
	})
}

// Refresh 處理換新權杖請求
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// 從 Header 讀取 Refresh Token (Authorization: Bearer <token>)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.App.ErrorJSON(w, errors.New("missing authorization header"), http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		h.App.ErrorJSON(w, errors.New("invalid authorization header"), http.StatusUnauthorized)
		return
	}
	refreshToken := parts[1]

	// 查詢 Refresh Token 是否存在
	var tokenRecord model.RefreshToken
	result := h.App.DB.Where("token = ?", refreshToken).First(&tokenRecord)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		h.App.ErrorJSON(w, errors.New("無效的換新權杖"), http.StatusUnauthorized)
		return
	}
	if result.Error != nil {
		h.App.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	// 檢查是否過期
	if tokenRecord.ExpiresAt.Before(time.Now()) {
		// 可以選擇在此刪除過期的 token
		h.App.DB.Delete(&tokenRecord)
		h.App.ErrorJSON(w, errors.New("換新權杖已過期"), http.StatusUnauthorized)
		return
	}

	// 查詢關聯的使用者
	var user model.User
	if err := h.App.DB.First(&user, tokenRecord.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.App.ErrorJSON(w, errors.New("user not found"), http.StatusUnauthorized)
		} else {
			h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		}
		return
	}

	// 產生新的 JWT token
	newToken, err := auth.GenerateToken(&user, time.Minute*time.Duration(h.App.Config.JWTExpiresIn), h.App.Config.JWTSecret)
	if err != nil {
		h.App.ErrorJSON(w, errors.New("無法產生新的認證 token"), http.StatusInternalServerError)
		return
	}

	// 刪除舊的 Refresh Token (Rotation 機制)
	if err := h.App.DB.Delete(&tokenRecord).Error; err != nil {
		h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// 產生並儲存新的 Refresh Token
	newRefreshToken, err := auth.CreateAndStoreRefreshToken(h.App.DB, user.ID, time.Minute*time.Duration(h.App.Config.JWTRefreshExpiresIn))
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法建立新的換新 token: %w", err), http.StatusInternalServerError)
		return
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":                  "權杖已更新",
		"token":                    newToken,
		"refresh_token":            newRefreshToken,
		"refresh_token_expires_in": h.App.Config.JWTRefreshExpiresIn * 60,
	})
}

// Me 處理取得個人資料請求
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// 從 Context 取得 user_id (由 AuthMiddleware 設定)
	userID, ok := r.Context().Value(middleware.UserContextKey).(uint)
	if !ok {
		h.App.ErrorJSON(w, errors.New("無法取得使用者資訊"), http.StatusUnauthorized)
		return
	}

	// 查詢使用者
	var user model.User
	if err := h.App.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.App.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		} else {
			h.App.ErrorJSON(w, err, http.StatusInternalServerError)
		}
		return
	}

	h.App.WriteJSON(w, http.StatusOK, apiHandler.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
