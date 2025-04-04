package handlers

import (
    "crs-backend/internal/database"
    "crs-backend/internal/models"
    "net/http"

    "github.com/gin-gonic/gin"
)

// رزرو بلیط
func BookTicket(c *gin.Context) {
    var input struct {
        FlightID uint `json:"flight_id"`
        UserID   uint `json:"user_id"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های نامعتبر"})
        return
    }

    var flight models.Flight
    if err := database.DB.First(&flight, input.FlightID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "پرواز یافت نشد"})
        return
    }

    if flight.Capacity <= 0 {
        c.JSON(http.StatusConflict, gin.H{"error": "ظرفیت پرواز تکمیل شده"})
        return
    }

    ticket := models.Ticket{
        FlightID: input.FlightID,
        UserID:   input.UserID,
        Status:   "Pending",
    }

    if err := database.DB.Create(&ticket).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در رزرو بلیط"})
        return
    }

    // کاهش ظرفیت پرواز
    database.DB.Model(&flight).Update("capacity", flight.Capacity-1)

    c.JSON(http.StatusCreated, gin.H{"message": "بلیط رزرو شد", "ticket": ticket})
}

// کنسل کردن بلیط
func CancelTicket(c *gin.Context) {
    id := c.Param("id")
    var ticket models.Ticket
    if err := database.DB.First(&ticket, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "بلیط یافت نشد"})
        return
    }

    database.DB.Delete(&ticket)
    c.JSON(http.StatusOK, gin.H{"message": "بلیط لغو شد"})
}
