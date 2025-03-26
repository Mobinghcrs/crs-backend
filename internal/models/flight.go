package models

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	ID           uint           `gorm:"primaryKey"`
	FlightNumber string         // شماره پرواز
	Origin       string         // مبدا
	Destination  string         // مقصد
	Departure    time.Time      // زمان پرواز از مبدا
	Arrival      time.Time      // زمان رسیدن به مقصد
	Price        float64        // قیمت بلیط
	Capacity     uint           // ظرفیت
	Tickets      []Ticket       `gorm:"foreignKey:FlightID"` // تعریف رابطه یک به چند
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

