package database

import (
	"log"
	"web-demo/config"
	"web-demo/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB 建立並回傳一個新的資料庫連線實例
func NewDB(cfg *config.AppConfig) (*gorm.DB, error) {
	// 建立資料庫連線
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Printf("資料庫已成功連線: %s", cfg.DBPath)

	// 自動遷移
	log.Println("正在執行資料庫遷移...")
	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.RefreshToken{})
	if err != nil {
		return nil, err
	}
	log.Println("資料庫遷移完成。")

	return db, nil
}

// CloseDB 用於安全地關閉指定的資料庫連線
func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("錯誤：無法從 GORM 取得底層 SQL 連線: %v", err)
		return
	}
	log.Println("正在關閉資料庫連線...")
	if err := sqlDB.Close(); err != nil {
		log.Printf("錯誤：關閉資料庫連線失敗: %v", err)
	}
}
