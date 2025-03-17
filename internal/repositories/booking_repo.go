package repositories

import (
	"crs-backend/internal/database"
	"crs-backend/internal/models"
)

// ایجاد رزرو جدید
func CreateBooking(booking *models.Booking) error {
	return database.DB.Create(booking).Error
}

// دریافت همه رزروها
func GetAllBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	err := database.DB.Preload("Ticket").Find(&bookings).Error
	return bookings, err
}

// دریافت رزرو بر اساس ID
func GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := database.DB.Preload("Ticket").First(&booking, id).Error
	return &booking, err
}

// لغو رزرو
func CancelBooking(id uint) error {
	return database.DB.Model(&models.Booking{}).Where("id = ?", id).Update("status", "canceled").Error
}
