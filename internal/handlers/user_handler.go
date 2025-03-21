package handlers

import (
	"crs-backend/internal/database"
	"crs-backend/internal/models"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// � دریافت لیست کاربران
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت لیست کاربران"})
		return
	}
	
	// پنهان کردن فیلدهای حساس
	for i := range users {
		users[i].PasswordHash = ""
	}
	
	c.JSON(http.StatusOK, users)
}

// � دریافت پروفایل کاربر
func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "کاربر احراز هویت نشده است"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "کاربر یافت نشد"})
		return
	}

	// پنهان کردن فیلدهای حساس
	user.PasswordHash = ""
	
	c.JSON(http.StatusOK, user)
}

// � دریافت اطلاعات کاربر خاص
func GetUser(c *gin.Context) {
	id := c.Param("id")
	
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "کاربر یافت نشد"})
		return
	}

	// پنهان کردن فیلدهای حساس
	user.PasswordHash = ""
	
	c.JSON(http.StatusOK, user)
}

// � حذف کاربر
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف کاربر"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "کاربر با موفقیت حذف شد"})
}
