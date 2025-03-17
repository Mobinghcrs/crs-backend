package repositories

import (
	"crs-backend/internal/database"
	"crs-backend/internal/models"
)

// ایجاد بلیط جدید
func CreateTicket(ticket *models.Ticket) error {
	return database.DB.Create(ticket).Error
}

// دریافت همه بلیط‌ها
func GetAllTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Find(&tickets).Error
	return tickets, err
}

// دریافت بلیط با ID
func GetTicketByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.First(&ticket, id).Error
	return &ticket, err
}

// به‌روزرسانی بلیط
func UpdateTicket(ticket *models.Ticket) error {
	return database.DB.Save(ticket).Error
}

// حذف بلیط
func DeleteTicket(id uint) error {
	return database.DB.Delete(&models.Ticket{}, id).Error
}
