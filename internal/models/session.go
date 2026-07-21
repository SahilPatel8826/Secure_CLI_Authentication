package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	UserID    uint
	Token     string `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
}
