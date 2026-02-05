package router

import (
	"errors"
	"net/http"
	"strings"
	"web-demo/handler/api"
	"web-demo/handler/web"
	"web-demo/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter 使用 chi 路由器並回傳一個 http.Handler 介面。
func NewRouter(app *server.Application) http.Handler {
	mux := chi.NewRouter()

	// 使用 chi 內建的中介軟體
	mux.Use(middleware.Logger)    // 記錄請求日誌
	mux.Use(middleware.Recoverer) // 從 panic 中恢復

	// 建立處理器實例
	apiHandler := api.NewAPIHandler(app)
	webHandler := web.NewWebHandler(app)
	authHandler := api.NewAuthHandler(app) // 新增 AuthHandler 實例

	// 註冊靜態檔案伺服器
	fileServer := http.FileServer(http.Dir("./public"))
	// chi 建議使用 Mount 來處理靜態檔案，並移除路徑前綴。
	mux.Mount("/static", http.StripPrefix("/static/", fileServer))

	// API 路由群組
	mux.Route("/api", func(r chi.Router) {
		r.Get("/users", apiHandler.GetAllUsers)
		r.Get("/users/{id}", apiHandler.GetUserByID)
		r.Get("/username/{name}", apiHandler.GetUserByName)

		// 認證相關路由
		r.Post("/auth/register", authHandler.Register) // 註冊路由
		r.Post("/auth/login", authHandler.Login)       // 登入路由
		r.Post("/auth/logout", authHandler.Logout)     // 登出路由
	})

	// 網頁路由
	mux.Get("/about", webHandler.AboutHandler)
	mux.Get("/", webHandler.HomeHandler)

	// 自訂 404 處理器
	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// 根據路徑前綴判斷是 API 請求還是網頁請求，並回傳對應的 404 錯誤。
		if strings.HasPrefix(r.URL.Path, "/api/") {
			app.ErrorJSON(w, errors.New("resource not found"), http.StatusNotFound)
		} else {
			webHandler.NotFoundHandler(w, r)
		}
	})

	return mux
}
