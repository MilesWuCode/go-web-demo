package routes

import "web-demo/handlers"

// apiRoutes 定義了所有 API 相關的路由
var apiRoutes = []Route{
	{
		Path:    "/api/echo",
		Handler: handlers.EchoHandler,
	},
}
