package database

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	"log"
	"os"

	"crs-backend/internal/models"
)

var DB *gorm.DB

func ConnectDB() {
	// بارگذاری متغیرهای محیطی از .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  خطا در بارگذاری فایل .env (ممکن است وجود نداشته باشد)")
	}

	// ساخت رشته اتصال
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "mobin2005G"),
		getEnv("DB_NAME", "crs_db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	// اتصال به دیتابیس
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ خطا در اتصال به دیتابیس: %v", err)
	}

	// نگه‌داشتن اتصال در متغیر DB
	DB = db

	// اجرای مهاجرت (ساخت جداول)
	migrateDatabase()

	fmt.Println("✅ اتصال به دیتابیس برقرار شد!")
}

// تابع مهاجرت برای ایجاد جداول
func migrateDatabase() {
	err := DB.AutoMigrate(&models.Ticket{}, &models.User{}, &models.Order{}, &models.Flight{})
	if err != nil {
		log.Fatalf("❌ خطا در اجرای مهاجرت دیتابیس: %v", err)
	}
	fmt.Println("✅ جداول دیتابیس ساخته شدند!")
}

// تابع کمکی برای دریافت متغیرهای محیطی با مقدار پیش‌فرض
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
