package utils

import (
	"time"
	"crs-backend/internal/models"


	"github.com/golang-jwt/jwt/v4"
)

// کلید مخفی برای امضای JWT
var jwtSecret = []byte("supersecretkey")

// ایجاد توکن JWT
func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // انقضای ۳ روزه
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// اعتبارسنجی توکن JWT
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
