package models

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	FlightNumber string    `json:"flight_number"` // شماره پرواز
	Origin       string    `json:"origin"`        // مبدأ
	Destination  string    `json:"destination"`   // مقصد
	Departure    time.Time `json:"departure"`     // زمان حرکت
	Arrival      time.Time `json:"arrival"`       // زمان رسیدن
	Price        float64   `json:"price"`         // قیمت
	Capacity     int       `json:"capacity"`      // ظرفیت
}
