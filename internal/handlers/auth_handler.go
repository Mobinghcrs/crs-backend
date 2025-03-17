package handlers

import (
	"crs-backend/internal/models"
	"crs-backend/internal/repositories"
	"crs-backend/internal/utils"
	
	"github.com/gin-gonic/gin"
	"net/http"
)

// ثبت‌نام کاربر جدید
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// هش کردن پسورد قبل از ذخیره
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در پردازش رمز عبور"})
		return
	}
	user.Password = hashedPassword

	// ایجاد کاربر در دیتابیس
	if err := repositories.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد کاربر"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "کاربر با موفقیت ثبت شد!"})
}

// ورود کاربر
func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// پیدا کردن کاربر در دیتابیس
	user, err := repositories.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ایمیل یا رمز عبور اشتباه است"})
		return
	}

	// بررسی رمز عبور
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ایمیل یا رمز عبور اشتباه است"})
		return
	}

	// تولید توکن JWT
	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد توکن"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
