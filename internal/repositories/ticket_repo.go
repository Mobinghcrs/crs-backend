package repositories

import (
	"crs-backend/internal/models"
	"gorm.io/gorm"
)

func CreateTicket(db *gorm.DB, ticket *models.Ticket) error {
	return db.Create(ticket).Error
}

func GetAllTickets(db *gorm.DB) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := db.Preload("Event").Find(&tickets).Error // اگر ارتباطی با Event وجود دارد
	return tickets, err
}

func GetTicketByID(db *gorm.DB, id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := db.Preload("Event").First(&ticket, id).Error
	return &ticket, err
}

func UpdateTicket(db *gorm.DB, ticket *models.Ticket) error {
	return db.Save(ticket).Error
}

func DeleteTicket(db *gorm.DB, id uint) error {
	return db.Delete(&models.Ticket{}, id).Error
}
