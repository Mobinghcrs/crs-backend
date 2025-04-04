package repositories

import (
	"booking-system/internal/models"
	"gorm.io/gorm"
	"strings"
	"fmt"
)

type FlightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) *FlightRepository {
	return &FlightRepository{db: db}
}

func (r *FlightRepository) CreateFlight(flight *models.Flight) error {
	result := r.db.Create(flight)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate") {
			return fmt.Errorf("flight number already exists")
		}
		return fmt.Errorf("database error: %v", result.Error)
	}
	return nil
}