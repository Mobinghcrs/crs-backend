package handlers

import (
	"net/http"
	"strconv"
    "errors"
	"booking-system/internal/models"
	"booking-system/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookingHandler struct {
	bookingService services.BookingService
	validator      *validator.Validate
}

func NewBookingHandler(bookingService services.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		validator:      validator.New(),
	}
}

type CreateBookingRequest struct {
	FlightID uint `json:"flight_id" binding:"required"`
	Seats    int  `json:"seats" binding:"required,min=1"`
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": h.formatValidationError(err)})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "دسترسی غیرمجاز"})
		return
	}

	booking := models.Booking{
		UserID:   userID.(uint),
		FlightID: req.FlightID,
		Seats:    req.Seats,
	}

	if err := h.bookingService.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": h.bookingErrorToMessage(err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "رزرو با موفقیت انجام شد",
		"id":      booking.ID,
	})
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "دسترسی غیرمجاز"})
		return
	}

	bookings, err := h.bookingService.GetUserBookings(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت رزروها"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	id := c.Param("id")
	bookingID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شناسه رزرو نامعتبر"})
		return
	}

	booking, err := h.bookingService.GetBookingByID(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "رزرو یافت نشد"})
		return
	}

	userID := c.MustGet("userID").(uint)
	if booking.UserID != userID && c.MustGet("role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "دسترسی غیرمجاز"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id := c.Param("id")
	bookingID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شناسه رزرو نامعتبر"})
		return
	}

	booking, err := h.bookingService.GetBookingByID(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "رزرو یافت نشد"})
		return
	}

	userID := c.MustGet("userID").(uint)
	if booking.UserID != userID && c.MustGet("role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "دسترسی غیرمجاز"})
		return
	}

	if err := h.bookingService.CancelBooking(uint(bookingID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در لغو رزرو"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "رزرو با موفقیت لغو شد"})
}

func (h *BookingHandler) formatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
	}
	return errors
}

func (h *BookingHandler) bookingErrorToMessage(err error) string {
	switch {
	case errors.Is(err, services.ErrFlightNotFound):
		return "پرواز مورد نظر یافت نشد"
	case errors.Is(err, services.ErrNotEnoughSeats):
		return "تعداد صندلی‌های درخواستی موجود نیست"
	case errors.Is(err, services.ErrDuplicateBooking):
		return "رزرو تکراری برای این پرواز وجود دارد"
	default:
		return "خطای داخلی سرور"
	}
}