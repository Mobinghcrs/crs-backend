package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Email    string `gorm:"unique;not null;type citext"` // استفاده از citext برای عدم حساسیت به حروف
    Password string `gorm:"not null"`
    Role     string `gorm:"not null;default:'user'"`
}