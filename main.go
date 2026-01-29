package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-demo/config"
	"web-demo/database"
	"web-demo/routes"
	"web-demo/server"
	"web-demo/view"
)

func main() {
	// --- 初始化 ---
	cfg := config.Get()

	// 初始化資料庫連線
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}
	defer database.CloseDB(db) // 使用 defer 確保程式結束時關閉資料庫

	// 初始化模板快取
	templateCache, err := view.NewCache()
	if err != nil {
		log.Fatalf("無法建立模板快取: %v", err)
	}

	// 建立 Application 實例，注入所有依賴
	app := server.Application{
		Config:        cfg,
		DB:            db,
		TemplateCache: templateCache,
	}

	router := routes.NewRouter(&app)

	// --- 設定與啟動伺服器 ---
	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		fmt.Printf("應用程式名稱: %s\n", cfg.AppName)
		fmt.Printf("伺服器已啟動: http://localhost%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("監聽失敗: %s\n", err)
		}
	}()

	// --- 設定優雅關閉 ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到關閉信號，伺服器準備關閉...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("伺服器關閉失敗: %v", err)
	}

	log.Println("伺服器已優雅地關閉")
}
