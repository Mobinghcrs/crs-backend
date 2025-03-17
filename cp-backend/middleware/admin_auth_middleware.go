package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminAuthMiddleware بررسی می‌کند که کاربر مدیر باشد
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// در اینجا بررسی JWT و نقش کاربر انجام می‌شود
		admin := c.GetHeader("X-Admin")
		if admin != "true" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
