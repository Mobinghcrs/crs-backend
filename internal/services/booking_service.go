package services

import (
	"booking-system/internal/models"
	"errors"
	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(booking *models.Booking) error
	GetBookingByID(id uint) (*models.Booking, error)
	CancelBooking(id uint) error
	GetUserBookings(userID uint) ([]models.Booking, error)
}

type bookingService struct {
	db *gorm.DB
}

func NewBookingService(db *gorm.DB) BookingService {
	return &bookingService{db: db}
}

func (s *bookingService) CreateBooking(booking *models.Booking) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var flight models.Flight
	if err := tx.First(&flight, booking.FlightID).Error; err != nil {
		tx.Rollback()
		return errors.New("پرواز یافت نشد")
	}

	if flight.AvailableSeats < booking.Seats {
		tx.Rollback()
		return errors.New("صندلی کافی موجود نیست")
	}

	if err := tx.Model(&flight).Update(
		"available_seats", 
		gorm.Expr("available_seats - ?", booking.Seats),
	).Error; err != nil {
		tx.Rollback()
		return errors.New("خطا در بروزرسانی صندلی‌ها")
	}

	booking.Status = "confirmed"
	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return errors.New("خطا در ایجاد رزرو")
	}

	tx.Commit()
	return nil
}

func (s *bookingService) GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := s.db.Preload("Flight").First(&booking, id).Error
	return &booking, err
}

func (s *bookingService) CancelBooking(id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var booking models.Booking
	if err := tx.First(&booking, id).Error; err != nil {
		tx.Rollback()
		return errors.New("رزرو یافت نشد")
	}

	if booking.Status == "cancelled" {
		tx.Rollback()
		return errors.New("رزرو قبلاً لغو شده است")
	}

	var flight models.Flight
	if err := tx.First(&flight, booking.FlightID).Error; err != nil {
		tx.Rollback()
		return errors.New("پرواز یافت نشد")
	}

	if err := tx.Model(&flight).Update(
		"available_seats", 
		gorm.Expr("available_seats + ?", booking.Seats),
	).Error; err != nil {
		tx.Rollback()
		return errors.New("خطا در بازگردانی صندلی‌ها")
	}

	if err := tx.Model(&booking).Update("status", "cancelled").Error; err != nil {
		tx.Rollback()
		return errors.New("خطا در لغو رزرو")
	}

	tx.Commit()
	return nil
}

func (s *bookingService) GetUserBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := s.db.Where("user_id = ?", userID).
		Preload("Flight").
		Find(&bookings).Error
	return bookings, err
}