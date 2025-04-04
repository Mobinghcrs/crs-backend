package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// بررسی نقش ادمین برای درخواست‌های کنترل پنل
func AdminAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        adminHeader := c.GetHeader("X-Admin")
        if adminHeader != "true" {
            c.JSON(http.StatusForbidden, gin.H{"error": "دسترسی غیرمجاز"})
            c.Abort()
            return
        }

        c.Next()
    }
}
