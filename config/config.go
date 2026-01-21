package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig 定義應用程式的所有設定參數
type AppConfig struct {
	AppName    string
	UploadPath string
	Port       string
}

// Load 初始化並回傳應用程式設定
func Load() *AppConfig {
	// 嘗試載入 .env 檔案 (非必須)
	if err := godotenv.Load(); err != nil {
		// 這裡改用 log.Println 提示即可，因為在正式環境可能只用系統環境變數
		log.Println("提示: 未讀取到 .env 檔案，將使用系統環境變數或預設值")
	}

	return &AppConfig{
		// 這裡設定您提到的固定參數
		AppName:    "Go Web Demo Application",
		UploadPath: "./uploads",

		// Port 仍然優先讀取環境變數，若無則使用預設值
		Port: getEnv("PORT", "3000"),
	}
}

// getEnv 讀取環境變數，若不存在則回傳預設值
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
