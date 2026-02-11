package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"web-demo/server"
	file "web-demo/util/file" // 導入新的 file 套件
)

// FileUploadHandler 結構用於處理檔案上傳相關請求
type FileUploadHandler struct {
	App *server.Application
}

// NewFileUploadHandler 建立並回傳一個新的 FileUploadHandler 實例
func NewFileUploadHandler(app *server.Application) *FileUploadHandler {
	return &FileUploadHandler{App: app}
}

// 處理檔案上傳
func (h *FileUploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// 限制上傳檔案大小為 10MB
	r.ParseMultipartForm(10 << 20) // 10MB

	fileContent, handler, err := r.FormFile("file") // "file" 是表單中檔案輸入欄位的名稱
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法獲取上傳檔案: %w", err), http.StatusBadRequest)
		return
	}
	defer fileContent.Close()

	// 產生一個混合英文和數字的唯一檔案名稱
	newFileName := file.GenerateEnglishMixedFileName(handler.Filename)
	destinationPath := filepath.Join("public", "upload", newFileName)

	// 確保 public/upload 目錄存在
	if _, err := os.Stat(filepath.Join("public", "upload")); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join("public", "upload"), 0755)
		if err != nil {
			h.App.ErrorJSON(w, fmt.Errorf("無法建立上傳目錄: %w", err), http.StatusInternalServerError)
			return
		}
	}

	dst, err := os.Create(destinationPath)
	if err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法建立目的檔案: %w", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 將上傳的檔案複製到目的檔案
	if _, err := io.Copy(dst, fileContent); err != nil {
		h.App.ErrorJSON(w, fmt.Errorf("無法儲存檔案: %w", err), http.StatusInternalServerError)
		return
	}

	h.App.WriteJSON(w, http.StatusOK, map[string]string{
		"message":   "檔案上傳成功",
		"file_name": newFileName,
		"file_path": "/upload/" + newFileName, // 提供相對於 public 的路徑
	})
}
