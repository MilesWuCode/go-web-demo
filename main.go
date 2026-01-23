package main

import (
	"fmt"
	"log"
	"net/http"
	"web-demo/config"
	"web-demo/database"
	"web-demo/handlers"
)

func main() {
	// 載入設定
	cfg := config.Load()

	// 初始化資料庫
	database.Init(cfg)

	// 設定靜態檔案伺服器
	// 將 /static/ 路徑的請求，交給 http.FileServer 處理
	// http.StripPrefix 會移除 URL 中的 /static/ 前綴，然後在 public/ 目錄中尋找檔案
	// 例如：/static/css/main.css -> public/css/main.css
	staticFileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// 註冊頁面路由
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/about", handlers.AboutHandler)

	// 註冊 API 路由
	http.HandleFunc("/api/echo", handlers.EchoHandler)

	// --- Server 啟動 ---
	addr := ":" + cfg.Port

	fmt.Printf("應用程式名稱: %s\n", cfg.AppName)
	fmt.Printf("伺服器已啟動: http://localhost%s\n", addr)

	// 啟動 Web Server
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
