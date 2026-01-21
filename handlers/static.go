package handlers

import "net/http"

// StaticHandler 回傳一個處理靜態檔案的 Handler
// 預設使用 "./static" 目錄
func StaticHandler() http.Handler {
	return http.FileServer(http.Dir("./static"))
}
