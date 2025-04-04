package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

const JWTSecret = "YwHrfYXHOv9eIdAIh0UMfM8s0+EYVpwpVuH7wZxd4jo=" // 🔑 جایگزین با کلید واقعی

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "booking-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		},
	)
}