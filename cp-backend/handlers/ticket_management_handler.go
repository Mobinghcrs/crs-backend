package handlers

import (
	"crs-backend/internal/models"
	"crs-backend/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// دریافت لیست بلیط‌ها
func ListTickets(c *gin.Context) {
	var tickets []models.Ticket
	result := database.DB.Find(&tickets)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت بلیط‌ها"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}

// ایجاد بلیط جدید
func CreateTicket(c *gin.Context) {
	var newTicket models.Ticket
	if err := c.ShouldBindJSON(&newTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "اطلاعات نادرست"})
		return
	}

	// ذخیره بلیط در دیتابیس
	result := database.DB.Create(&newTicket)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد بلیط"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "بلیط با موفقیت ایجاد شد", "ticket": newTicket})
}

// ویرایش اطلاعات بلیط
func UpdateTicket(c *gin.Context) {
	var ticket models.Ticket
	id := c.Param("id")

	// جستجوی بلیط در دیتابیس
	result := database.DB.First(&ticket, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "بلیط یافت نشد"})
		return
	}

	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "اطلاعات نادرست"})
		return
	}

	// ذخیره تغییرات
	database.DB.Save(&ticket)
	c.JSON(http.StatusOK, gin.H{"message": "اطلاعات بلیط بروزرسانی شد", "ticket": ticket})
}

// حذف بلیط
func DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Ticket{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف بلیط"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "بلیط با موفقیت حذف شد"})
}
