package models

import (
	"time"
	"gorm.io/gorm"
)

type Ticket struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Price     float64        `json:"price"`
	Available uint           `json:"available"`
	Departure time.Time      `json:"departure"`
	Arrival   time.Time      `json:"arrival"`
	Status    string `gorm:"size:50;default:'pending'"`

	FlightID  uint   `json:"flight_id"`
	// رابطه به مدل هواپیما (Flight)
	Flight    Flight `json:"flight"`

	UserID    uint   `json:"user_id"`  // کلید خارجی برای کاربر
	OrderID   uint   `json:"order_id"` // کلید خارجی برای سفارش
	  
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
