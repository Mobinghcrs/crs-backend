package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB متغیر گلوبال برای پایگاه داده
var DB *gorm.DB

// Config ساختار تنظیمات پایگاه داده
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection تابع برای اتصال به پایگاه داده
func NewConnection(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}

// Migrate اجرای مهاجرت‌های پایگاه داده
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		// اینجا مدل‌های پایگاه داده را اضافه کنید
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Database migrated successfully!")
}
