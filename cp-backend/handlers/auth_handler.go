package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// Register برای ثبت‌نام کاربر
func Register(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "ثبت‌نام موفق"})
}

// Login برای ورود کاربر
func Login(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "ورود موفق"})
}
