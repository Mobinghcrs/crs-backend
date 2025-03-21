package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"localhost",   // جایگزینی با مقادیر واقعی
		"postgres",    // کاربر دیتابیس
		"mobin2005G", // پسورد
		"crs_db",      // نام دیتابیس
		"5432",        // پورت
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db
	return DB
}
