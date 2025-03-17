package handlers

import (
	"crs-backend/internal/models"
	"crs-backend/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// دریافت لیست کاربران
func ListUsers(c *gin.Context) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت کاربران"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ایجاد کاربر جدید
func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "اطلاعات نادرست"})
		return
	}

	// ذخیره کاربر در دیتابیس
	result := database.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد کاربر"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "کاربر با موفقیت ایجاد شد", "user": newUser})
}

// ویرایش اطلاعات کاربر
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// جستجوی کاربر در دیتابیس
	result := database.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "کاربر یافت نشد"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "اطلاعات نادرست"})
		return
	}

	// ذخیره تغییرات
	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "اطلاعات کاربر بروزرسانی شد", "user": user})
}

// حذف کاربر
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف کاربر"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "کاربر با موفقیت حذف شد"})
}
