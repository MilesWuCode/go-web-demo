package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// AppConfig 定義應用程式的所有設定參數
type AppConfig struct {
	AppName    string
	UploadPath string
	Port       string
	DBPath     string
}

var (
	instance *AppConfig
	once     sync.Once
)

// Get 使用單例模式回傳 AppConfig 的唯一實例。
// 設定檔的載入邏輯只會在第一次被呼叫時執行。
func Get() *AppConfig {
	// once.Do 會保證在多執行緒環境下，內部的函式也只會被執行一次。
	once.Do(func() {
		// 嘗試載入 .env 檔案 (非必須)
		if err := godotenv.Load(); err != nil {
			log.Println("提示: 未讀取到 .env 檔案，將使用系統環境變數或預設值")
		}

		// 載入設定並將其賦值給套件級別的 instance 變數
		instance = &AppConfig{
			AppName:    "Go Web Demo Application",
			UploadPath: "./uploads",
			DBPath:     "db.sqlite",
			Port:       getEnv("PORT", "3000"),
		}
	})
	return instance
}

// getEnv 讀取環境變數，若不存在則回傳預設值
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
