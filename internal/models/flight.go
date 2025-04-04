package models

import (
    "time"
    "gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	FlightNumber   string    `gorm:"uniqueIndex;size:20" json:"flight_number"`
	Origin         string    `gorm:"size:100" json:"origin"`
	Destination    string    `gorm:"size:100" json:"destination"`
	DepartureTime  time.Time `json:"departure_time"`
	Capacity       int       `json:"capacity"`
	AvailableSeats int       `json:"available_seats"`
	BasePrice      float64   `json:"base_price"`
	Airline        string    `gorm:"size:50" json:"airline"`
	Aircraft       string    `gorm:"size:50" json:"aircraft"`
}