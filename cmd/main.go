package main

import (
    "crs-backend/internal/models"
    "crs-backend/internal/repositories"
    "crs-backend/routes"
    "log"
    "os"
    "time"
    "net/http"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/gin-gonic/gin"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("❌ خطا در بارگذاری فایل .env: ", err)
    }

    dbConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    }

    db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), dbConfig)
    if err != nil {
        log.Fatal("❌ خطای اتصال به دیتابیس: ", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("❌ خطای دریافت connection pool: ", err)
    }
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("❌ خطای ping به دیتابیس: ", err)
    }

    migrate(db)

    // 🔴 ایجاد نمونه ریپازیتوری
	eventRepo := repositories.NewEventRepository(db)
    
    // 🔴 انتقال وابستگی‌ها به روتر
    router := routes.SetupRouter(db, eventRepo)

    startServer(router)
}

func migrate(db *gorm.DB) {
	// فعال‌سازی قابلیت UUID در PostgreSQL
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// مهاجرت مدل‌ها با ترتیب صحیح وابستگی
	models := []interface{}{
		&models.User{},
		&models.Event{},
		&models.Ticket{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("❌ خطای مهاجرت برای مدل %T: %v", model, err)
		}
	}
	log.Println("✅ مهاجرت دیتابیس با موفقیت انجام شد")
}

func startServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("� سرور در حال اجرا روی پورت %s...", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("❌ خطای اجرای سرور: ", err)
	}
}
