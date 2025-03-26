package models

import "time"

type SecurityLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Action    string    `gorm:"size:100;not null"`
	IPAddress string    `gorm:"size:45;not null"` // IPv6 پشتیبانی
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
