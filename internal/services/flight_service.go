package services

import (
	"booking-system/internal/models"
	"errors"
	"gorm.io/gorm"
)

// اینترفیس سرویس
type FlightService interface {
	CreateFlight(flight *models.Flight) error
	UpdateFlight(flight *models.Flight) error
	DeleteFlight(id uint) error
	GetFlights(page, limit int) ([]models.Flight, error)
}

// پیاده‌سازی سرویس (توجه به نام ساختار)
type flightService struct {
	db *gorm.DB
}

// سازنده سرویس
func NewFlightService(db *gorm.DB) FlightService {
	return &flightService{db: db}
}

// CreateFlight ایجاد پرواز جدید
func (s *flightService) CreateFlight(flight *models.Flight) error {
	if flight.AvailableSeats > flight.Capacity {
		return errors.New("available seats exceed capacity")
	}
	return s.db.Create(flight).Error
}

// UpdateFlight به‌روزرسانی پرواز
func (s *flightService) UpdateFlight(flight *models.Flight) error {
	return s.db.Save(flight).Error
}

// DeleteFlight حذف پرواز
func (s *flightService) DeleteFlight(id uint) error {
	return s.db.Delete(&models.Flight{}, id).Error
}

// GetFlights دریافت لیست پروازها
func (s *flightService) GetFlights(page, limit int) ([]models.Flight, error) {
	var flights []models.Flight
	offset := (page - 1) * limit
	return flights, s.db.Offset(offset).Limit(limit).Find(&flights).Error
}