package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ReadIntQuery 是一個輔助函式，用於從 URL 查詢字串中讀取一個整數。
// 如果查詢參數不存在或不是一個有效的整數，則回傳預設值。
func (app *Application) ReadIntQuery(r *http.Request, key string, defaultValue int) int {
	strValue := r.URL.Query().Get(key)
	if strValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	if intValue < 1 {
		return defaultValue
	}

	return intValue
}

// WriteJSON 是一個輔助函式，用於方便地回傳 JSON 格式的回應。
// 它會自動設定 Content-Type 標頭並處理 JSON 編碼。
func (app *Application) WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("錯誤：編碼 JSON 失敗: %v", err)
	}
}

// ErrorJSON 是一個輔助函式，用於回傳一個標準格式的 JSON 錯誤。
func (app *Application) ErrorJSON(w http.ResponseWriter, err error, status int) {
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	app.WriteJSON(w, status, errorResponse)
}

// ParseRequestForm 是一個輔助函式，用於解析來自請求主體的表單資料。
// 它會自動解析 application/x-www-form-urlencoded 和 multipart/form-data。
func (app *Application) ParseRequestForm(w http.ResponseWriter, r *http.Request) error {
	// Default max memory for ParseMultipartForm is 32MB. Can be configured if needed.
	// ParseForm handles both application/x-www-form-urlencoded and multipart/form-data
	// for POST, PUT, and PATCH requests. For GET, it parses the URL query string.
	err := r.ParseForm()
	if err != nil {
		return err
	}
	return nil
}
