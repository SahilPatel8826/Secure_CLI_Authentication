package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	// IsLocked            bool       `gorm:"default:false"`
	MFAEnabled          bool       `gorm:"default:false"`
	MFASecret           string     `gorm:"default:null"`
	FailedLoginAttempts int        `gorm:"default:0"`
	LockedUntil         *time.Time `gorm:"default:null"`
	LastLogin           *time.Time `gorm:"default:null"`
}
