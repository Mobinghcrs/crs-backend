package handlers

import (
    "crs-backend/internal/database"
    "crs-backend/internal/models"
    "net/http"

    "github.com/gin-gonic/gin"
)

// ایجاد سفارش
func CreateOrder(c *gin.Context) {
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "داده‌های نامعتبر"})
        return
    }

    if err := database.DB.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ایجاد سفارش"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "سفارش ایجاد شد", "order": order})
}

// دریافت سفارشات کاربر
func GetOrders(c *gin.Context) {
    var orders []models.Order
    database.DB.Find(&orders)
    c.JSON(http.StatusOK, orders)
}
