package repositories

import (
	"crs-backend/internal/models"
	"gorm.io/gorm"
)

func CreateBooking(db *gorm.DB, booking *models.Booking) error {
	return db.Create(booking).Error
}

func GetAllBookings(db *gorm.DB) ([]models.Booking, error) {
	var bookings []models.Booking
	err := db.Preload("Ticket").Find(&bookings).Error // اضافه شدن Preload برای ارتباط
	return bookings, err
}

func GetBookingByID(db *gorm.DB, id uint) (*models.Booking, error) {
	var booking models.Booking
	err := db.Preload("Ticket").First(&booking, id).Error // اضافه شدن Preload
	return &booking, err
}

func CancelBooking(db *gorm.DB, id uint) error {
	return db.Model(&models.Booking{}).Where("id = ?", id).Update("status", "canceled").Error
}
func UpdateBookingStatus(db *gorm.DB, id uint, status string) error {
	return db.Model(&models.Booking{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// افزودن تابع DeleteBooking
func DeleteBooking(db *gorm.DB, id uint) error {
	return db.Delete(&models.Booking{}, id).Error
}
