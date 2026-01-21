package main

import (
	"fmt"
	"log"
	"net/http"
	"web-demo/handlers"
)

func main() {
	// 設定靜態檔案伺服器
	fs := http.FileServer(http.Dir("./static"))

	// 註冊路由
	// API 路由
	http.HandleFunc("/api/echo", handlers.EchoHandler)
	
	// 靜態檔案路由 (根路徑匹配)
	http.Handle("/", fs)

	port := ":3000"
	fmt.Printf("伺服器已啟動: http://localhost%s\n", port)
	fmt.Printf("API 測試路徑: http://localhost%s/api/echo\n", port)

	// 啟動 Web Server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}