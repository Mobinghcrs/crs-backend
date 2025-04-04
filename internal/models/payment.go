package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	BookingID   uint   `json:"booking_id"`
	Amount      int64  `json:"amount"`
	Status      string `json:"status"` // pending, success, failed
	TransactionID string `json:"transaction_id"`
}
