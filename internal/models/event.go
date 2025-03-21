package models

import (
    "time"
    "gorm.io/gorm"
)

type Event struct {
    gorm.Model
    Title       string     `gorm:"size:255;not null" validate:"required"`
    Description string     `gorm:"type:text" validate:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"` 
	EndDate     time.Time `json:"end_date" binding:"required"`
    EndTime     time.Time  `gorm:"not null"`
    Location    string     `gorm:"size:255;not null"`
    CreatorID   uint       `gorm:"not null"`
    Tickets     []Ticket   `gorm:"foreignKey:EventID"`
}