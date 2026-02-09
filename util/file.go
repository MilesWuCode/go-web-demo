package util

import (
	"fmt"
	"math/rand"
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
