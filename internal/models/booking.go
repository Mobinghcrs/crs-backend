package models

import (
	"gorm.io/gorm"
	"time"
)

type BookingStatus string

const (
	BookingPending  BookingStatus = "pending"
	BookingPaid     BookingStatus = "paid"
	BookingCanceled BookingStatus = "canceled"
)

type Booking struct {
	gorm.Model
	UserID     uint          `gorm:"not null;index"`         // ارتباط با کاربر
	User       User          `gorm:"foreignKey:UserID"`      // رابطه belongsTo
	FlightID   uint          `gorm:"not null;index"`         // ارتباط با پرواز
	Flight     Flight        `gorm:"foreignKey:FlightID"`    // رابطه belongsTo
	Seats      int           `gorm:"not null;min=1"`         // حداقل 1 صندلی
	Status     BookingStatus `gorm:"type:varchar(20);index"` // استفاده از enum
	TotalPrice float64       `gorm:"not null"`               // قیمت کل
	PaymentID  string        `gorm:"type:varchar(100)"`      // شناسه پرداخت
	BookedAt   time.Time     `gorm:"default:CURRENT_TIMESTAMP"` // زمان رزرو
}

// جدول پرداخت جداگانه برای مدیریت تراکنش‌ها
type Payment struct {
	gorm.Model
	BookingID    uint    `gorm:"not null;uniqueIndex"`
	Amount       float64 `gorm:"not null"`
	Currency     string  `gorm:"type:varchar(3);not null"`
	PaymentMethod string `gorm:"type:varchar(50);not null"`
	TransactionID string `gorm:"type:varchar(100);uniqueIndex"`
}