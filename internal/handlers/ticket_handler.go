package handlers

import (
	"crs-backend/internal/models"
	"crs-backend/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ایجاد بلیط جدید
func CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repositories.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد بلیط"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// دریافت همه بلیط‌ها
func GetAllTickets(c *gin.Context) {
	tickets, err := repositories.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت بلیط‌ها"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// دریافت بلیط بر اساس ID
func GetTicketByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ticket, err := repositories.GetTicketByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "بلیط پیدا نشد"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// به‌روزرسانی بلیط
func UpdateTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ticket, err := repositories.GetTicketByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "بلیط پیدا نشد"})
		return
	}

	if err := c.ShouldBindJSON(ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repositories.UpdateTicket(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در به‌روزرسانی بلیط"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// حذف بلیط
func DeleteTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repositories.DeleteTicket(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف بلیط"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "بلیط حذف شد"})
}
func AddTicket(c *gin.Context) {
	// اینجا منطق اضافه کردن بلیط جدید رو اضافه کن
	c.JSON(http.StatusOK, gin.H{"message": "بلیط اضافه شد"})
}