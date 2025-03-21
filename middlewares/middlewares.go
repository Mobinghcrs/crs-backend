package middlewares 

import (
	"net/http"
	"os" 
	"strings"
	"time"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// تنظیمات JWT
var (
	// اصلاح شده: خواندن کلید از محیط اجرا
	secretKey     = []byte(os.Getenv("JWT_SECRET")) 
	// اصلاح شده: پیکربندی مدت زمان از محیط اجرا
	tokenDuration = 24 * time.Hour
)

// ساختار کلیمزهای توکن
type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ------------------------------------------------------------
// میدلورهای پایه
// ------------------------------------------------------------

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}

// ------------------------------------------------------------
// سیستم احراز هویت JWT
// ------------------------------------------------------------

func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", jwt.ErrTokenRequiredClaimMissing
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", jwt.ErrTokenInvalidId
	}

	return tokenParts[1], nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil {
			sendError(c, http.StatusUnauthorized, "TOKEN_REQUIRED", "نیاز به توکن احراز هویت")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				sendError(c, http.StatusUnauthorized, "TOKEN_EXPIRED", "توکن منقضی شده است")
				return
			}
			sendError(c, http.StatusUnauthorized, "INVALID_TOKEN", "توکن نامعتبر")
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			sendError(c, http.StatusUnauthorized, "INVALID_CLAIMS", "ساختار توکن نامعتبر")
			return
		}

		// ذخیره اطلاعات کاربر در Context
		c.Set("authUserID", claims.UserID)
		c.Set("authUserRole", claims.Role)
		c.Next()
	}
}

// ------------------------------------------------------------
// توابع کمکی
// ------------------------------------------------------------

func sendError(c *gin.Context, code int, errorCode string, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":      errorCode,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339Nano),
		"path":      c.Request.URL.Path,
	})
}

func GenerateJWT(userID uint, role string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "crs-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
