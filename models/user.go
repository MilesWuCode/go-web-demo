package models

import "gorm.io/gorm"

// User 模型對應到資料庫中的 'users' 資料表
type User struct {
	gorm.Model        // 包含了 ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"unique"`
	Email      string `gorm:"unique"`
	Password   string // 儲存密碼 (注意：應儲存雜湊值而非明文)
	Posts      []Post // 一對多關聯：一個使用者可以有多個 Post
}