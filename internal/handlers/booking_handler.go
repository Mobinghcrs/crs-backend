package handlers

import (
	"net/http"
	"strconv"

	"crs-backend/internal/models"
	"crs-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingHandler struct {
	db *gorm.DB
}

// تابع سازنده جدید برای هندلر
func NewBookingHandler(db *gorm.DB) *BookingHandler {
	return &BookingHandler{db: db}
}

// ایجاد رزرو جدید
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_BOOKING_DATA",
			"message": "داده‌های رزرو نامعتبر است",
		})
		return
	}

	if err := repositories.CreateBooking(h.db, &booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "BOOKING_CREATION_FAILED",
			"message": "خطا در ایجاد رزرو",
		})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// دریافت تمامی رزروها
func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	bookings, err := repositories.GetAllBookings(h.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FETCH_BOOKINGS_FAILED",
			"message": "خطا در دریافت لیست رزروها",
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// دریافت رزرو بر اساس ID
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	booking, err := repositories.GetBookingByID(h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "BOOKING_NOT_FOUND",
			"message": "رزرو مورد نظر یافت نشد",
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// بروزرسانی وضعیت رزرو
func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_STATUS_DATA",
			"message": "وضعیت ارسال شده نامعتبر است",
		})
		return
	}

	if err := repositories.UpdateBookingStatus(h.db, uint(id), request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "BOOKING_STATUS_UPDATE_FAILED",
			"message": "خطا در بروزرسانی وضعیت رزرو",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "وضعیت رزرو با موفقیت بروزرسانی شد",
	})
}

// حذف رزرو
func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := repositories.DeleteBooking(h.db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "BOOKING_DELETION_FAILED",
			"message": "خطا در حذف رزرو",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "رزرو با موفقیت حذف شد",
	})
}
