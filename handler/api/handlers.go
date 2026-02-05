package api

import "web-demo/server"

// UserResponse 用於定義回傳給前端的使用者資料結構，以隱藏密碼等敏感資訊
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type APIHandler struct {
	App *server.Application
}

func NewAPIHandler(app *server.Application) *APIHandler {
	return &APIHandler{App: app}
}
