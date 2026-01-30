package web

import (
	"bytes"
	"log"
	"net/http"
)

// Render 是一個輔助函式，用來方便地渲染模板
func (h *WebHandler) Render(w http.ResponseWriter, status int, page string) {
	// 從快取中取得模板
	ts, ok := h.App.TemplateCache[page]
	if !ok {
		log.Printf("錯誤：模板 %s 不存在", page)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 使用一個緩衝區來執行模板，以便在出錯時能捕捉錯誤
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, nil)
	if err != nil {
		log.Printf("錯誤：渲染模板失敗: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 寫入 HTTP 狀態碼和模板內容
	w.WriteHeader(status)
	buf.WriteTo(w)
}
