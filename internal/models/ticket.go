package models

import (
	"gorm.io/gorm"
	"time"
)

// مدل بلیط
type Ticket struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	EventID     uint      	   `gorm:"index"`
	Price       float64        `gorm:"not null" json:"price"`
	Available   int            `gorm:"not null" json:"available"` // تعداد بلیط‌های موجود
	Departure   time.Time      `json:"departure"`                 // زمان حرکت
	Arrival     time.Time      `json:"arrival"`                   // زمان رسیدن
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // حذف نرم (Soft Delete)
}
