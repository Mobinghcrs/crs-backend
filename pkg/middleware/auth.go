package middleware

import (
    "context"
    "net/http"
    "strings"
    "crs-backend/internal/models"
    "gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, `{"error": "احراز هویت ضروری است"}`, http.StatusUnauthorized)
                return
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")
            
            var user models.User
            if err := db.Where("token = ?", token).First(&user).Error; err != nil {
                http.Error(w, `{"error": "توکن نامعتبر"}`, http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), "user", &user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func RoleMiddleware(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := r.Context().Value("user").(*models.User)
            
            if user.Role != role {
                http.Error(w, `{"error": "دسترسی غیرمجاز"}`, http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
