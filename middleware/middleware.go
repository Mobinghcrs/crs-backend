package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// کلید مخفی برای JWT
var secretKey = []byte("your-secret-key")

// AuthMiddleware: بررسی احراز هویت کاربران با JWT
func AuthMiddleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		// دریافت توکن از هدر Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن ارائه نشده است"})
			c.Abort()
			return
		}

		// حذف "Bearer " از ابتدای توکن
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// بررسی و پردازش توکن JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		// بررسی اعتبار توکن
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن نامعتبر است"})
			c.Abort()
			return
		}

		// ادامه پردازش درخواست
		c.Next()
	}
}

// AdminMiddleware: بررسی نقش ادمین برای مدیریت
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// دریافت توکن از هدر Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن ارائه نشده است"})
			c.Abort()
			return
		}

		// حذف "Bearer " از ابتدای توکن
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// بررسی و پردازش توکن JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		// بررسی اعتبار توکن
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن نامعتبر است"})
			c.Abort()
			return
		}

		// بررسی نقش کاربر (ادمین یا نه)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["role"] != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "دسترسی غیرمجاز"})
			c.Abort()
			return
		}

		// ادامه پردازش درخواست
		c.Next()
	}
}
