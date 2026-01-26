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
	// 1. 載入設定
	cfg := config.Get()

	// 2. 資料庫
	database.Init(cfg)

	// 3. 模板快取
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

	// 6. 設定與啟動伺服器
	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 在一個新的 goroutine 中啟動伺服器
	go func() {
		fmt.Printf("應用程式名稱: %s\n", cfg.AppName)
		fmt.Printf("伺服器已啟動: http://localhost%s\n", addr)
		// ListenAndServe 會阻塞，直到發生錯誤
		// 當我們呼叫 Shutdown() 時，它會回傳 http.ErrServerClosed，這不是一個真正的錯誤
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("監聽失敗: %s\n", err)
		}
	}()

	// --- 設定優雅關閉 ---
	// 建立一個 channel 來接收作業系統的信號
	quit := make(chan os.Signal, 1)
	// 我們希望捕捉 SIGINT (Ctrl+C) 和 SIGTERM (kill 指令)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 程式會在這裡阻塞，直到收到上述信號
	<-quit
	log.Println("收到關閉信號，伺服器準備關閉...")

	// 建立一個有超時的 context，給予 5 秒的時間讓現有請求處理完畢
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 呼叫 Shutdown()，它會優雅地關閉伺服器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("伺服器關閉失敗: %v", err)
	}

	// 在伺服器成功關閉後，關閉資料庫連線
	database.CloseDB()

	log.Println("伺服器已優雅地關閉")
}
