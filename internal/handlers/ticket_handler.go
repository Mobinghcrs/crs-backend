package handlers

import (
	"crs-backend/internal/models"
	"crs-backend/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TicketHandler struct {
	db *gorm.DB
}

func NewTicketHandler(db *gorm.DB) *TicketHandler {
	return &TicketHandler{db: db}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_TICKET_DATA",
			"message": "داده‌های بلیط نامعتبر است",
		})
		return
	}

	if err := repositories.CreateTicket(h.db, &ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "TICKET_CREATION_FAILED",
			"message": "خطا در ایجاد بلیط",
		})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	// اصلاح: اضافه کردن h.db به عنوان پارامتر اول
	ticket, err := repositories.GetTicketByID(h.db, uint(id))
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "TICKET_NOT_FOUND",
			"message": "بلیط مورد نظر یافت نشد",
		})
		return
	}
	
	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	tickets, err := repositories.GetAllTickets(h.db) // اضافه کردن h.db
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "TICKETS_FETCH_FAILED",
			"message": "خطا در دریافت لیست بلیط‌ها",
		})
		return
	}
	
	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_TICKET_DATA",
			"message": "داده‌های بلیط نامعتبر است",
		})
		return
	}

	// اصلاح: اضافه کردن h.db به عنوان پارامتر اول
	if err := repositories.UpdateTicket(h.db, &ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "TICKET_UPDATE_FAILED",
			"message": "خطا در بروزرسانی بلیط",
		})
		return
	}
	
	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	// اصلاح: اضافه کردن h.db به عنوان پارامتر اول
	if err := repositories.DeleteTicket(h.db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "TICKET_DELETE_FAILED",
			"message": "خطا در حذف بلیط",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "بلیط با موفقیت حذف شد",
	})
}
