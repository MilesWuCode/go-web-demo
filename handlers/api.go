package handlers

import (
	"encoding/json"
	"net/http"
)

// EchoHandler 處理 /api/echo 請求
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{"data": "hello world"}
	json.NewEncoder(w).Encode(data)
}
