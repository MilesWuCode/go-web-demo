package models

import "gorm.io/gorm"

// Post 模型對應到資料庫中的 'posts' 資料表
type Post struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint // 外鍵，對應到 User 的 ID
	User    User // 關聯：一個 Post 從屬於一個 User
}