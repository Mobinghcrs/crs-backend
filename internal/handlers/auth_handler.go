package handlers

import (
	"crs-backend/internal/database"
	"crs-backend/internal/models"
	"crs-backend/internal/repositories"
	
	"net/http"
	
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)
func Login1(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های ورودی نامعتبر"})
		return
	}

	// دریافت کاربر از دیتابیس
	user, err := repositories.GetUserByUsername(database.DB, credentials.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "نام کاربری یا رمز عبور اشتباه"})
		return
	}
	if user.Username == "" {
		user.Username = user.Email // یا تولید خودکار
	}
	// بررسی تطابق رمز عبور
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(credentials.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "نام کاربری یا رمز عبور اشتباه"})
		return
	}

	// TODO: ایجاد توکن JWT
	c.JSON(http.StatusOK, gin.H{"message": "ورود موفقیت آمیز"})
}
// 📌 ثبت نام کاربر جدید
func Register1(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های ورودی نامعتبر"})
		return
	}

	// بررسی وجود کاربر با نام کاربری تکراری
	exists, err := repositories.UsernameExists(database.DB, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در بررسی نام کاربری"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "نام کاربری قبلا ثبت شده است"})
		return
	}

	// هش کردن رمز عبور
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در پردازش رمز عبور"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	// ایجاد کاربر در دیتابیس
	if err := repositories.CreateUser(database.DB, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد کاربر"})
		return
	}

	// پنهان کردن فیلدهای حساس
	user.Password = ""
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "کاربر با موفقیت ثبت شد",
		"user":    user,
	})
}
