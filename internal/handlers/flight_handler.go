package handlers

import (
	"net/http"
	"strconv"
	//"fmt"
	"strings"
	"booking-system/internal/models"
	"booking-system/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FlightHandler struct {
	flightService services.FlightService
	validator     *validator.Validate
}

func NewFlightHandler(flightService services.FlightService) *FlightHandler {
	return &FlightHandler{
		flightService: flightService,
		validator:     validator.New(),
	}
}

// CreateFlight ایجاد پرواز جدید
func (h *FlightHandler) CreateFlight(c *gin.Context) {
	var flight models.Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": h.formatValidationError(err)})
		return
	}

	// اعتبارسنجی داده‌ها
	if err := h.validator.Struct(flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": h.formatValidationError(err)})
		return
	}

	if err := h.flightService.CreateFlight(&flight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": h.handleDatabaseError(err)})
		return
	}

	c.JSON(http.StatusCreated, flight)
}

// GetFlights دریافت لیست پروازها
func (h *FlightHandler) GetFlights(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	flights, err := h.flightService.GetFlights(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت لیست پروازها"})
		return
	}

	c.JSON(http.StatusOK, flights)
}

// UpdateFlight به‌روزرسانی پرواز
func (h *FlightHandler) UpdateFlight(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شناسه پرواز نامعتبر"})
		return
	}

	var flight models.Flight
	if err := c.ShouldBindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های ورودی نامعتبر"})
		return
	}

	flight.ID = uint(flightID)
	if err := h.flightService.UpdateFlight(&flight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در به‌روزرسانی پرواز"})
		return
	}

	c.JSON(http.StatusOK, flight)
}

// DeleteFlight حذف پرواز
func (h *FlightHandler) DeleteFlight(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شناسه پرواز نامعتبر"})
		return
	}

	if err := h.flightService.DeleteFlight(uint(flightID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در حذف پرواز"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "پرواز با موفقیت حذف شد"})
}

// formatValidationError مدیریت خطاهای اعتبارسنجی
func (h *FlightHandler) formatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
	}
	return errors
}

// handleDatabaseError مدیریت خطاهای دیتابیس
func (h *FlightHandler) handleDatabaseError(err error) string {
	if strings.Contains(err.Error(), "duplicate") {
		return "اطلاعات تکراری"
	}
	return "خطای داخلی سرور"
}