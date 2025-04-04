package models

import (
	"time"
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	TicketID  uint      `gorm:"not null" json:"ticket_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	UserPhone string    `gorm:"not null" json:"user_phone"`
	UserEmail string    `gorm:"not null" json:"user_email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Ticket    Ticket    `gorm:"foreignKey:TicketID" json:"ticket"`
}
