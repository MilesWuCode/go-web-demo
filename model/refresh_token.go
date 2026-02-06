package model

import (
	"time"
)

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"uniqueIndex;size:255"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
