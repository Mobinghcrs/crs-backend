package handlers

import (
	"crs-backend/internal/database"
	"crs-backend/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 📌 متد ایجاد پرداخت و هدایت به درگاه زرین‌پال
func CreatePayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ورودی نامعتبر است"})
		return
	}

	// ثبت پرداخت در دیتابیس با وضعیت "pending"
	payment.Status = "pending"
	database.DB.Create(&payment)

	// ایجاد لینک پرداخت (مثال برای زرین‌پال)
	paymentURL := fmt.Sprintf("https://www.zarinpal.com/pg/StartPay/%d", payment.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":     "لینک پرداخت ایجاد شد",
		"payment_url": paymentURL,
	})
}

// 📌 متد بررسی وضعیت پرداخت
func VerifyPayment(c *gin.Context) {
	transactionID := c.Query("transaction_id") // دریافت شناسه پرداخت از URL

	var payment models.Payment
	if err := database.DB.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "پرداخت یافت نشد"})
		return
	}

	// فرض کنیم که زرین‌پال یه API برای تأیید پرداخت داره و وضعیت رو می‌گیریم
	// اینجا فقط یه شبیه‌سازی هست
	isPaid := true // در حالت واقعی، از API زرین‌پال بررسی می‌کنیم

	if isPaid {
		payment.Status = "success"
	} else {
		payment.Status = "failed"
	}

	database.DB.Save(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "وضعیت پرداخت بروزرسانی شد", "status": payment.Status})
}
