package utils

import (
    "time"
	"fmt"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("eVYo2eGcPhD+S5P5AJ8SvM5hlS6fjxWQsqEC0vPq3mM=")

// ✅ تغییر نام تابع به VerifyToken و Export کردن
func VerifyToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
}

func GenerateJWT(userID uint, role string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })
    return token.SignedString(jwtSecret)
}
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("الگوریتم ناشناخته: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}