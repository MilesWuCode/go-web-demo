package model

import (
	"gorm.io/gorm"
)

// SocialProvider 模型用於儲存第三方登入的用戶資訊
type SocialProvider struct {
	gorm.Model
	UserID     uint   `gorm:"index"`                     // 關聯到內部 User 模型
	Provider   string `gorm:"size:50;not null"`          // 登入提供者 (e.g., "google")
	ProviderID string `gorm:"size:255;not null;unique"`  // 提供者提供的唯一 ID
	Email      string `gorm:"size:255"`                  // 提供者提供的 Email
	Name       string `gorm:"size:255"`                  // 提供者提供的名稱
	// 其他從提供者獲取的資訊可以根據需求添加
}
