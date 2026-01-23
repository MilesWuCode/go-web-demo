package main

import (
	"fmt"
	"log"
	"net/http"
	"web-demo/config"
	"web-demo/database"
	"web-demo/routes"
)

func main() {
	// 1. 載入設定
	cfg := config.Load()

	// 2. 初始化資料庫
	database.Init(cfg)

	// 3. 建立並註冊所有路由
	router := routes.NewRouter()

	// 4. 啟動伺服器
	addr := ":" + cfg.Port
	fmt.Printf("應用程式名稱: %s\n", cfg.AppName)
	fmt.Printf("伺服器已啟動: http://localhost%s\n", addr)

	// 將我們自己建立的 router 傳遞給伺服器
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
