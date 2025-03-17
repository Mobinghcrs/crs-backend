package models

import (
	"time"

	"gorm.io/gorm"
)

// مدل رزرو بلیط
type Booking struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TicketID  uint           `gorm:"not null" json:"ticket_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Quantity  int            `gorm:"not null" json:"quantity"` // تعداد بلیط‌های رزرو شده
	Status    string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // وضعیت رزرو (pending, confirmed, canceled)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // حذف نرم

	// ارتباط با مدل Ticket
	Ticket Ticket `gorm:"foreignKey:TicketID" json:"ticket"`
}
