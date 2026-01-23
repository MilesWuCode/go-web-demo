package main

import (
	"fmt"
	"log"
	"net/http"
	"web-demo/config"
	"web-demo/database"
	"web-demo/routes"
	"web-demo/server"
	"web-demo/view"
)

func main() {
	// 1. 載入設定
	cfg := config.Get()

	// 2. 初始化資料庫
	database.Init(cfg)

	// 3. 初始化模板快取
	templateCache, err := view.NewCache()
	if err != nil {
		log.Fatalf("無法建立模板快取: %v", err)
	}

	// 4. 建立 Application 實例，注入依賴
	app := server.Application{
		TemplateCache: templateCache,
	}

	// 5. 建立並註冊所有路由
	router := routes.NewRouter(&app)

	// 6. 啟動伺服器
	addr := ":" + cfg.Port
	fmt.Printf("應用程式名稱: %s\n", cfg.AppName)
	fmt.Printf("伺服器已啟動: http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
