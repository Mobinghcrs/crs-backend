package models

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// مدل کاربر
type User struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Username  string `gorm:"size:100;uniqueIndex;not null"`
    PasswordHash string         `gorm:"not null" json:"-"`
	Password  string `gorm:"size:255;not null" json:"-"` //
    FullName     string         `gorm:"not null" json:"full_name"`
    Email     string `gorm:"size:255;uniqueIndex;not null"`
    Role      string `gorm:"size:50;default:'user'"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	TwoFactorEnabled  bool   `gorm:"default:false"`
  	TwoFactorSecret   string 
  	RecoveryCodes     string `gorm:"type:text"` // JSON array
	Tickets  []Ticket `gorm:"foreignKey:UserID"` // مشخص کردن کلید خارجی برای رابطه یک به چند

}

func UserRoleMiddleware(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != requiredRole {
            c.JSON(http.StatusForbidden, gin.H{"error": "شما دسترسی لازم را ندارید"})
            c.Abort()
            return
        }

        c.Next()
    }
}
