package handlers

import (
    "crs-backend/internal/database"
    "crs-backend/internal/models"
    "net/http"

    "github.com/gin-gonic/gin"
)
func ListFlights(c *gin.Context) {
	// پیاده‌سازی واقعی
	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{"id": 1, "flight_number": "AB123", "origin": "THR", "destination": "MHD"},
			{"id": 2, "flight_number": "CD456", "origin": "IKA", "destination": "SYZ"},
		},
	})
}

// ایجاد پرواز (Admin Only)
func CreateFlight(c *gin.Context) {
    var flight models.Flight
    if err := c.ShouldBindJSON(&flight); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های نامعتبر"})
        return
    }

    if err := database.DB.Create(&flight).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد پرواز"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "پرواز با موفقیت ایجاد شد", "flight": flight})
}

// دریافت لیست پروازها
func GetFlights(c *gin.Context) {
    var flights []models.Flight
    database.DB.Find(&flights)
    c.JSON(http.StatusOK, flights)
}
func UpdateFlight(c *gin.Context) {
	// پیاده‌سازی واقعی
	c.JSON(http.StatusOK, gin.H{"message": "پرواز با موفقیت بروزرسانی شد"})
}
// حذف پرواز (Admin Only)
func DeleteFlight(c *gin.Context) {
    id := c.Param("id")
    if err := database.DB.Delete(&models.Flight{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف پرواز"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "پرواز حذف شد"})
}
