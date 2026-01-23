package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// NotFoundHandler 處理所有未找到的頁面請求
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// 設定 HTTP 狀態碼為 404
	w.WriteHeader(http.StatusNotFound)

	fp := filepath.Join("templates", "404.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		// 如果 404 模板本身有問題，就回傳一個簡單的文字錯誤
		http.Error(w, "404 Not Found", http.StatusNotFound)
		log.Printf("錯誤：解析 404 模板失敗 %v", err)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("錯誤：執行 404 模板失敗 %v", err)
	}
}

// HomeHandler 處理根路徑 "/" 的請求，並作為所有未匹配路由的「捕獲」處理器
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// 檢查是否為未匹配的 API 請求
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r) // 回傳純文字的 404 Not Found
		return
	}

	// 檢查請求的路徑是否為根路徑，如果不是，則顯示 HTML 404 頁面
	if r.URL.Path != "/" {
		NotFoundHandler(w, r)
		return
	}

	// --- 如果是根路徑，則渲染首頁模板 ---
	fp := filepath.Join("templates", "index.html")

	// 解析模板
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("錯誤：解析模板失敗 %v", err)
		return
	}

	// 執行模板
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("錯誤：執行模板失敗 %v", err)
	}
}

// AboutHandler 處理 "/about" 的請求，渲染 about.html 模板
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "about.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("錯誤：解析模板失敗 %v", err)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("錯誤：執行模板失敗 %v", err)
	}
}
