package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 設定靜態檔案目錄
	fs := http.FileServer(http.Dir("./static"))

	// API 路由：/api/echo
	http.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := map[string]string{"data": "hello world"}
		json.NewEncoder(w).Encode(data)
	})

	// 將根路徑 "/" 的請求交給 file server 處理
	// 注意：http.Handle 會匹配所有路徑，所以 API 路由必須先定義或使用更精確的匹配
	http.Handle("/", fs)

	port := ":3000"
	fmt.Printf("伺服器即將啟動，請瀏覽: http://localhost%s\n", port)

	// 啟動 Web Server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
