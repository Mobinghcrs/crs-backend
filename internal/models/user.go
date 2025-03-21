package models

import (
	"time"

	"gorm.io/gorm"
)

// مدل کاربر
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username string          `gorm:"type:text;not null;default:''"`
	PasswordHash string 	 `gorm:"type:text;not null;default:'default_value'" json:"password_hash"`     
	FullName  string         `gorm:"not null" json:"full_name"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"type:varchar(20);default:'user'" json:"role"` // نقش: user, admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
