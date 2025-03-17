package repositories

import (
	"crs-backend/internal/database"
	"crs-backend/internal/models"
)

// ایجاد کاربر جدید
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// دریافت کاربر بر اساس ایمیل
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
