package server

import (
	"encoding/json"
	"log"
	"net/http"
)

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
