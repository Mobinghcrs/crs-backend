package models

import (
	"time"
	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID     uint    `gorm:"not null" json:"user_id"`
	TotalPrice float64 `gorm:"not null" json:"total_price"`
	Status     string  `gorm:"type:varchar(20);default:'pending'" json:"status"` // وضعیت سفارش

	// رابطه یک به چند با بلیط‌ها
	Tickets []Ticket `gorm:"foreignKey:OrderID" json:"tickets"`

	// رابطه هر سفارش به یک کاربر
	User User `gorm:"foreignKey:UserID" json:"user"`
}
