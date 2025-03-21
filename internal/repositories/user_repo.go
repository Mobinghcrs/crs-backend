package repositories

import (
	"crs-backend/internal/models"
	"gorm.io/gorm"
)

// 📌 ایجاد کاربر جدید
func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

// 📌 بررسی وجود نام کاربری
func UsernameExists(db *gorm.DB, username string) (bool, error) {
	var count int64
	err := db.Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).
		Error
	return count > 0, err
}

// 📌 دریافت کاربر با ایمیل
func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// 📌 دریافت کاربر با نام کاربری
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	return &user, err
}
