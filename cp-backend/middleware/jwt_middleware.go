package middleware

import (
	"crs-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "توکن احراز هویت یافت نشد"})
			return
		}

		token, err := utils.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "توکن نامعتبر"})
			return
		}

		// 🔴 اصلاح: استخراج صحیح claims از توکن
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "فرمت توکن نامعتبر"})
			return
		}

		// 🔴 اصلاح: دسترسی به claimها با کلیدهای صحیح
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}
