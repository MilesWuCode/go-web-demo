package routes

import (
	"net/http"
	"web-demo/server"
)

// NewRouter 現在接收 Application 實例，並註冊其上的方法作為 Handler
func NewRouter(app *server.Application) *http.ServeMux {
	mux := http.NewServeMux()

	// 註冊靜態檔案伺服器
	dir := http.Dir("./public")
	fileServer := http.FileServer(dir)
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// 註冊 API 路由
	mux.HandleFunc("/api/echo", app.EchoHandler)
	mux.HandleFunc("/api/users", app.GetAllUsers)

	// 註冊頁面路由
	mux.HandleFunc("/about", app.AboutHandler)
	mux.HandleFunc("/", app.HomeHandler) // HomeHandler 作為捕獲所有路由的處理器

	return mux
}
