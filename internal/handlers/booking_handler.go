package handlers

import (
	"crs-backend/internal/models"
	"crs-backend/internal/notifications"
	"crs-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"strconv"
)

// ایجاد رزرو جدید
func CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// بررسی ظرفیت بلیط قبل از رزرو
	ticket, err := repositories.GetTicketByID(booking.TicketID)
	if err != nil || ticket.Available < booking.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ظرفیت کافی نیست!"})
		return
	}

	// ثبت رزرو در دیتابیس
	if err := repositories.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد رزرو"})
		return
	}

	// کاهش ظرفیت بلیط پس از رزرو موفق
	ticket.Available -= booking.Quantity
	repositories.UpdateTicket(ticket)

	c.JSON(http.StatusCreated, booking)
		// اطلاعات کاربر (باید از توکن JWT استخراج بشه)
		userEmail := "user@example.com"
		userPhone := "09123456789"
	
		// 📩 ارسال پیامک تأیید رزرو
		smsMessage := "رزرو شما با موفقیت انجام شد! جزئیات در ایمیل شما ارسال شد."
		err = notifications.SendSMS(userPhone, smsMessage)  
		if err != nil {
			fmt.Println("❌ خطا در ارسال پیامک:", err)
		}
	
		// 📧 ارسال ایمیل تأیید رزرو
		emailSubject := "تأیید رزرو بلیط شما"
		emailBody := "رزرو شما با موفقیت ثبت شد. لطفاً جزئیات را بررسی کنید."
		err = notifications.SendEmail(userEmail, emailSubject, emailBody)
		if err != nil {
			fmt.Println("❌ خطا در ارسال ایمیل:", err)
		}
	
		// پاسخ نهایی به کاربر
		c.JSON(http.StatusOK, gin.H{"message": "رزرو با موفقیت انجام شد و اعلان‌ها ارسال شدند."})
	
}

// دریافت همه رزروها
func GetAllBookings(c *gin.Context) {
	bookings, err := repositories.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت رزروها"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// دریافت رزرو بر اساس ID
func GetBookingByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	booking, err := repositories.GetBookingByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "رزرو پیدا نشد"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// لغو رزرو
func CancelBooking(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repositories.CancelBooking(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در لغو رزرو"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "رزرو لغو شد"})
}
func GetAvailableTickets(c *gin.Context) {
	// اینجا منطق دریافت بلیط‌ها رو اضافه کن
	c.JSON(http.StatusOK, gin.H{"message": "لیست بلیط‌های موجود"})
}

// دریافت لیست رزروهای کاربر
func GetUserBookings(c *gin.Context) {
	// اینجا منطق دریافت رزروهای کاربر رو اضافه کن
	c.JSON(http.StatusOK, gin.H{"message": "لیست رزروهای کاربر"})
}