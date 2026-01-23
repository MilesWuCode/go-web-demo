package routes

import (
	"net/http"
	"web-demo/handlers"
)

// Route 定義了一個路由的結構
type Route struct {
	Path    string
	Handler http.HandlerFunc
}

// webRoutes 定義了所有與頁面渲染相關的路由
var webRoutes = []Route{
	{
		Path:    "/about",
		Handler: handlers.AboutHandler,
	},
	{
		// 將根路徑 "/" (捕獲所有未匹配頁面) 放在最後
		Path:    "/",
		Handler: handlers.HomeHandler,
	},
}
