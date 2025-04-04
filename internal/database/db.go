package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"crs-backend/internal/models"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func ConnectDB() {
    // بارگذاری .env (حذف پیام خطا اگر فایل وجود نداشت)
    _ = godotenv.Load()

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_USER", "postgres"),
        getEnv("DB_PASSWORD", "mobin2005G"), 
        getEnv("DB_NAME", "crs_db"),
        getEnv("DB_PORT", "5432"),
        getEnv("DB_SSLMODE", "disable"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("خطا در اتصال به دیتابیس: %v", err)
    }

    // اجرای مهاجرت با ترتیب صحیح
    if err := DB.AutoMigrate(
        &models.User{},
        &models.Flight{},
        &models.Order{},
        &models.Ticket{},
    ); err != nil {
        log.Fatalf("خطا در مهاجرت دیتابیس: %v", err)
    }
}


// توجه کنید که ترتیب مهاجرت اهمیت دارد؛ ابتدا مدل‌هایی که توسط سایر مدل‌ها ارجاع داده می‌شوند (User, Order, Flight) و سپس Ticket
func migrateDatabase() {
	err := DB.AutoMigrate( &models.User{}, &models.Order{}, &models.Flight{}, &models.Ticket{}, )
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

func NewConnection(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
