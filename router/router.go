package router

import (
	"errors"
	"net/http"
	"strings"
	authHandler "web-demo/handler/api/auth"      // Import the new auth handler package
	googleVerifyHandler "web-demo/handler/api/auth" // Import the new google verify handler package
	fileHandler "web-demo/handler/api/file"      // Import the new file handler package
	s3FileHandler "web-demo/handler/api/s3_file" // Import the new s3_file handler package
	userHandler "web-demo/handler/api/user"      // Import the new user handler package
	"web-demo/handler/web"
	"web-demo/middleware"
	"web-demo/server"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter 使用 chi 路由器並回傳一個 http.Handler 介面。
func NewRouter(app *server.Application) http.Handler {
	mux := chi.NewRouter()

	// 使用 chi 內建的中介軟體
	mux.Use(chiMiddleware.Logger)    // 記錄請求日誌
	mux.Use(chiMiddleware.Recoverer) // 從 panic 中恢復

	// 建立處理器實例
	webHandler := web.NewWebHandler(app)
	authH := authHandler.NewAuthHandler(app)
	googleVerifyH := googleVerifyHandler.NewGoogleVerifyHandler(app)
	userH := userHandler.NewUserHandler(app)
	fileUploadH := fileHandler.NewFileUploadHandler(app)
	s3FileUploadHandler := s3FileHandler.NewS3FileUploadHandler(app) // New S3 file upload handler instance

	// 註冊靜態檔案伺服器
	fileServer := http.FileServer(http.Dir("./public"))
	// chi 建議使用 Mount 來處理靜態檔案，並移除路徑前綴。
	mux.Mount("/static", http.StripPrefix("/static/", fileServer))

	// API 路由群組
	mux.Route("/api", func(r chi.Router) {
		r.Get("/users", userH.GetAllUsers)
		r.Get("/users/{id}", userH.GetUserByID)
		r.Get("/username/{name}", userH.GetUserByName)

		// 檔案上傳路由
		r.Post("/demo-upload-file", fileUploadH.UploadFile)
		r.Post("/demo-s3-file-upload", s3FileUploadHandler.UploadS3File)

		// 認證相關路由
		r.Post("/auth/register", authH.Register)     // 註冊路由
		r.Post("/auth/login", authH.Login)           // 登入路由
		r.Post("/auth/refresh-token", authH.Refresh) // 換新權杖路由
		r.Post("/auth/google/verify-token", googleVerifyH.VerifyGoogleIDToken) // Google ID Token 驗證路由

		// 登出路由，應用 JWT 驗證中間件
		r.With(middleware.AuthMiddleware(app)).Post("/auth/logout", authH.Logout)
		r.With(middleware.AuthMiddleware(app)).Get("/auth/me", authH.Me)
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
