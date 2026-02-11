package util

import (
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

// GenerateEnglishMixedFileName 產生一個混合英文和數字的檔案名稱
func GenerateEnglishMixedFileName(originalFileName string) string {
	ext := filepath.Ext(originalFileName)

	// 使用時間戳和隨機字串確保唯一性
	timestamp := time.Now().UnixNano()
	randomString := GenerateRandomString(8) // 產生 8 個字元的隨機字串

	return fmt.Sprintf("%d_%s%s", timestamp, randomString, ext)
}

// GenerateRandomString 產生指定長度的隨機英數字元串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// DetectContentType 偵測檔案的 Content-Type
func DetectContentType(file io.ReadSeeker, filename string) (string, error) {
	// 讀取檔案開頭部分以偵測 Content-Type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("無法讀取檔案內容以偵測類型: %w", err)
	}
	// 將讀取指針重置回檔案開頭
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("無法重置檔案讀取指針: %w", err)
	}

	// 優先使用 http.DetectContentType 偵測 MIME 類型
	contentType := http.DetectContentType(buffer)
	if contentType == "application/octet-stream" {
		// 如果 http.DetectContentType 仍是通用類型，則嘗試使用副檔名判斷
		extType := mime.TypeByExtension(filepath.Ext(filename))
		if extType != "" {
			contentType = extType
		}
	}
	return contentType, nil
}
