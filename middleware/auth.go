package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey 是一個私有類型，用於在 context 中儲存鍵。
type contextKey string

// UserContextKey 是用於在請求 context 中儲存使用者 ID 的鍵。
const UserContextKey = contextKey("userID")

// AppContext 定義了中介層所需的應用程式依賴。
type AppContext interface {
	ErrorJSON(w http.ResponseWriter, err error, status int)
	GetJWTSecret() string
}

// AuthMiddleware 是一個驗證 JWT token 的中間件。
// 它會從 Authorization 標頭中提取 token，驗證它，並將使用者 ID 儲存到請求 context 中。
func AuthMiddleware(app AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 從 Authorization 標頭中獲取 token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.ErrorJSON(w, errors.New("缺少認證 token"), http.StatusUnauthorized)
				return
			}

			// 檢查 token 格式是否為 "Bearer <token>"
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				app.ErrorJSON(w, errors.New("認證 token 格式錯誤"), http.StatusUnauthorized)
				return
			}

			tokenString := headerParts[1]

			// 解析和驗證 token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 檢查簽名方法
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("無效的簽名方法")
				}
				return []byte(app.GetJWTSecret()), nil
			})

			if err != nil {
				app.ErrorJSON(w, errors.New("無效或過期的 token: "+err.Error()), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				app.ErrorJSON(w, errors.New("無效的 token"), http.StatusUnauthorized)
				return
			}

			// 從 token 中提取聲明
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				app.ErrorJSON(w, errors.New("無法解析 token 聲明"), http.StatusUnauthorized)
				return
			}

			// 獲取使用者 ID
			userID, ok := claims["user_id"].(float64) // JWT claims數字通常解析為 float64
			if !ok {
				app.ErrorJSON(w, errors.New("token 中缺少使用者 ID"), http.StatusUnauthorized)
				return
			}

			// 將使用者 ID 儲存到請求 context 中
			ctx := context.WithValue(r.Context(), UserContextKey, uint(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
