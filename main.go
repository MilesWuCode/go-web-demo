package main

import (
	"fmt"
	"log"
	"net/http"
	"web-demo/config"
	"web-demo/handlers"
)

func main() {
	// 載入設定
	cfg := config.Load()

	// 註冊路由
	// API 路由
	http.HandleFunc("/api/echo", handlers.EchoHandler)

	// 靜態檔案路由 (根路徑匹配)
	http.Handle("/", handlers.StaticHandler())

	// 確保 port 有冒號前綴
	addr := ":" + cfg.Port

	fmt.Printf("應用程式名稱: %s\n", cfg.AppName)
	fmt.Printf("檔案上傳路徑: %s\n", cfg.UploadPath)
	fmt.Printf("伺服器已啟動: http://localhost%s\n", addr)
	fmt.Printf("API 測試路徑: http://localhost%s/api/echo\n", addr)

	// 啟動 Web Server
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}