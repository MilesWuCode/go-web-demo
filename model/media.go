package model

import (
	"gorm.io/gorm"
)

// Media 模型對應到資料庫中的 'media' 資料表
type Media struct {
	gorm.Model     // 包含了 ID, CreatedAt, UpdatedAt, DeletedAt
	ModelName      string
	CollectionName string
	FileName       string
}
