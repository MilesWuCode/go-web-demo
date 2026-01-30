package routes

import (
	"net/http"
	"web-demo/handlers/api"
	"web-demo/handlers/web"
	"web-demo/server"
)

// NewRouter 現在接收 Application 實例，並註冊其上的方法作為 Handler
func NewRouter(app *server.Application) *http.ServeMux {
	mux := http.NewServeMux()

	// 建立 handler 實例
	apiHandler := api.NewAPIHandler(app)
	webHandler := web.NewWebHandler(app)

	// 註冊靜態檔案伺服器
	dir := http.Dir("./public")
	fileServer := http.FileServer(dir)
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// 註冊 API 路由
	mux.HandleFunc("/api/echo", apiHandler.EchoHandler)
	mux.HandleFunc("GET /api/users", apiHandler.GetAllUsers)
	mux.HandleFunc("GET /api/users/{id}", apiHandler.GetUserByID)

	// 註冊頁面路由
	mux.HandleFunc("/about", webHandler.AboutHandler)
	mux.HandleFunc("/", webHandler.HomeHandler) // HomeHandler 作為捕獲所有路由的處理器

	return mux
}
