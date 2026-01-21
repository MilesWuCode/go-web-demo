package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"web-demo/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// 載入 .env 檔案
	// 如果沒有 .env 檔案，這裡不會報錯，直接使用系統環境變數或預設值
	if err := godotenv.Load(); err != nil {
		log.Println("未發現 .env 檔案，將使用系統環境變數")
	}

	// 註冊路由
	// API 路由
	http.HandleFunc("/api/echo", handlers.EchoHandler)

	// 靜態檔案路由 (根路徑匹配)
	http.Handle("/", handlers.StaticHandler())

	// 從環境變數取得 PORT，預設為 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 確保 port 有冒號前綴
	addr := ":" + port

	fmt.Printf("伺服器已啟動: http://localhost%s\n", addr)
	fmt.Printf("API 測試路徑: http://localhost%s/api/echo\n", addr)

	// 啟動 Web Server
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}