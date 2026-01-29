package database

import (
	"log"
	"web-demo/config"
	"web-demo/models" // 引入 models 套件

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化資料庫連線並執行自動遷移
func Init(cfg *config.AppConfig) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}
	log.Printf("資料庫已成功連線: %s", cfg.DBPath)

	// 自動遷移
	log.Println("正在執行資料庫遷移...")

	err = DB.AutoMigrate(&models.User{}, &models.Post{}) // 加入 Post 模型

	if err != nil {
		log.Fatalf("資料庫遷移失敗: %v", err)
	}

	log.Println("資料庫遷移完成。")
}

// CloseDB 用於安全地關閉資料庫連線
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("錯誤：無法從 GORM 取得底層 SQL 連線: %v", err)
		return
	}
	log.Println("正在關閉資料庫連線...")
	if err := sqlDB.Close(); err != nil {
		log.Printf("錯誤：關閉資料庫連線失敗: %v", err)
	}
}
