package repositories

import (
	"crs-backend/internal/models"
	"gorm.io/gorm"
)

type IEventRepository interface {
	Create(event *models.Event) error
	GetAllEvents() ([]models.Event, error)
	GetEventByID(id uint) (*models.Event, error)
	Update(event *models.Event) error
	Delete(id uint) error
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// همه‌ی متدها به صورت receiver functions
func (r *EventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *EventRepository) GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Preload("Tickets").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) GetEventByID(id uint) (*models.Event, error) {
	var event models.Event
	if err := r.db.Preload("Tickets").First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *EventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}
