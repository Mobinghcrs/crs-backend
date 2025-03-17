package middleware

import (
	"crs-backend/internal/utils"
	
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// بررسی توکن JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن ارائه نشده است"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		_, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن نامعتبر است"})
			c.Abort()
			return
		}

		c.Next()
	}
}
