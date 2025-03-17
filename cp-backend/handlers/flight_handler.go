package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"crs-backend/internal/database"
	"crs-backend/internal/models"
)

// ListFlights دریافت لیست تمام پروازها (اختیاری)
func ListFlights(c *gin.Context) {
	var flights []models.Flight
	if err := database.DB.Find(&flights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت لیست پروازها"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"flights": flights})
}

// CreateFlight ایجاد یا افزودن پرواز جدید
func CreateFlight(c *gin.Context) {
	var flight models.Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ورودی نامعتبر"})
		return
	}
	if err := database.DB.Create(&flight).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد پرواز"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "پرواز با موفقیت ایجاد شد", "flight": flight})
}

// UpdateFlight ویرایش اطلاعات پرواز
func UpdateFlight(c *gin.Context) {
	id := c.Param("id")
	var flight models.Flight
	if err := database.DB.First(&flight, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "پرواز یافت نشد"})
		return
	}
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ورودی نامعتبر"})
		return
	}
	if err := database.DB.Save(&flight).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در به‌روزرسانی پرواز"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "پرواز به‌روزرسانی شد", "flight": flight})
}

// DeleteFlight حذف پرواز
func DeleteFlight(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Flight{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف پرواز"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "پرواز حذف شد"})
}
