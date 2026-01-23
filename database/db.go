package database

import (
	"log"
	"web-demo/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化資料庫連線
func Init(cfg *config.AppConfig) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}

	log.Printf("資料庫已成功連線: %s", cfg.DBPath)
}
