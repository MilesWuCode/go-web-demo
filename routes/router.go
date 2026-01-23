package routes

import "net/http"

// NewRouter 建立並回傳一個設定好所有應用程式路由的 http.ServeMux
func NewRouter() *http.ServeMux {
	// 建立一個新的 ServeMux，避免使用全域的 DefaultServeMux
	mux := http.NewServeMux()

	// 註冊所有 Web 路由
	for _, route := range webRoutes {
		mux.HandleFunc(route.Path, route.Handler)
	}

	// 註冊所有 API 路由
	for _, route := range apiRoutes {
		mux.HandleFunc(route.Path, route.Handler)
	}

	// 處理靜態檔案
	// 靜態檔案的 Handler 需要另外註冊，因為它不是 http.HandlerFunc 類型
	dir := http.Dir("./public")
	fileServer := http.FileServer(dir)
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}
