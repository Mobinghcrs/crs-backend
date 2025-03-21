package handlers

import (
	"crs-backend/internal/models"
	"crs-backend/internal/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	repo repositories.IEventRepository // فقط از ریپازیتوری استفاده می‌کنیم
}

func NewEventHandler(repo repositories.IEventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

// GetAllEvents - دریافت همه رویدادها
func (h *EventHandler) GetAllEvents(c *gin.Context) { 
	events, err := h.repo.GetAllEvents() // ✅ فراخوانی از طریق ریپازیتوری
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "خطا در دریافت رویدادها",
		})
		return
	}

	if len(events) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":    []interface{}{},
			"message": "هیچ رویدادی یافت نشد",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": events})
}

// CreateEvent - ایجاد رویداد جدید
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ // ✅ استفاده از استاندارد HTTP
			"code":    "INVALID_INPUT",
			"message": "داده‌های ورودی نامعتبر",
		})
		return
	}

	// اعتبارسنجی پیشرفته
	if event.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "MISSING_TITLE",
			"message": "عنوان رویداد الزامی است",
		})
		return
	}

	if event.StartDate.IsZero() || event.EndDate.Before(event.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_DATES",
			"message": "تاریخ‌های شروع و پایان نامعتبر",
		})
		return
	}

	// تنظیم مقادیر پیش‌فرض
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	if err := h.repo.Create(&event); err != nil { // ✅ استفاده از ریپازیتوری
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "EVENT_CREATION_FAILED",
			"message": "خطا در ایجاد رویداد",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    event,
		"message": "رویداد با موفقیت ایجاد شد",
	})
}

// UpdateEvent - به‌روزرسانی رویداد
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_INPUT",
			"message": "داده‌های ورودی نامعتبر",
		})
		return
	}

	// دریافت رویداد موجود
	existingEvent, err := h.repo.GetEventByID(uint(id)) // ✅ استفاده از ریپازیتوری
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "EVENT_NOT_FOUND",
			"message": "رویداد مورد نظر یافت نشد",
		})
		return
	}

	// به‌روزرسانی فیلدها
	existingEvent.Title = event.Title
	existingEvent.Description = event.Description
	existingEvent.StartDate = event.StartDate
	existingEvent.EndDate = event.EndDate
	existingEvent.UpdatedAt = time.Now()

	if err := h.repo.Update(existingEvent); err != nil { // ✅ استفاده از ریپازیتوری
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "UPDATE_FAILED",
			"message": "خطا در به‌روزرسانی رویداد",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    existingEvent,
		"message": "رویداد با موفقیت به‌روزرسانی شد",
	})
}

// DeleteEvent - حذف رویداد
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// بررسی وجود رویداد
	if _, err := h.repo.GetEventByID(uint(id)); err != nil { // ✅ استفاده از ریپازیتوری
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "EVENT_NOT_FOUND",
			"message": "رویداد مورد نظر یافت نشد",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil { // ✅ استفاده از ریپازیتوری
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DELETE_FAILED",
			"message": "خطا در حذف رویداد",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "رویداد با موفقیت حذف شد",
	})
}

// GetEventByID - دریافت رویداد بر اساس ID
func (h *EventHandler) GetEventByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	event, err := h.repo.GetEventByID(uint(id)) // ✅ استفاده از ریپازیتوری
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "EVENT_NOT_FOUND",
			"message": "رویداد مورد نظر یافت نشد",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": event})
}
