package main

import (
	"log"
	"booking-system/internal/handlers"
	"booking-system/internal/models"
	"booking-system/internal/services"
	"booking-system/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. اتصال به دیتابیس
	dsn := "host=localhost user=postgres password=mobin2005G dbname=booking port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("خطا در اتصال به دیتابیس: ", err)
	}

	// 2. تنظیمات Connection Pool
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 3. اجرای مهاجرت‌ها
	if err := db.AutoMigrate(
		&models.User{},
		&models.Flight{},
		&models.Booking{},
	); err != nil {
		log.Fatal("خطا در مهاجرت دیتابیس: ", err)
	}

	// 4. ایجاد روتر Gin
	r := gin.Default()

	// 5. میدلورهای عمومی
	r.Use(
		middleware.CORSMiddleware(),
		gin.Recovery(),
	)

	// 6. مقداردهی سرویس‌ها
	authService := services.NewAuthService(db)
	flightService := services.NewFlightService(db)
	bookingService := services.NewBookingService(db)

	// 7. مقداردهی هندلرها
	authHandler := handlers.NewAuthHandler(authService)
	flightHandler := handlers.NewFlightHandler(flightService)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// 8. Routeهای عمومی
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/signup", authHandler.SignUp)
	r.POST("/login", authHandler.Login)
	r.GET("/flights", flightHandler.GetFlights) // دسترسی عمومی

	// 9. گروه Routeهای احراز هویت شده
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/bookings", bookingHandler.CreateBooking)
		authGroup.GET("/bookings", bookingHandler.GetUserBookings)
	}

	// 10. گروه Routeهای مدیریتی (ادمین)
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AdminMiddleware())
	{
		// مدیریت پروازها
		adminGroup.POST("/flights", flightHandler.CreateFlight)
		adminGroup.PUT("/flights/:id", flightHandler.UpdateFlight)
		adminGroup.DELETE("/flights/:id", flightHandler.DeleteFlight)

		// مدیریت کاربران
		adminGroup.GET("/users", authHandler.ListUsers)
	}

	// 11. اجرای سرور
	log.Println("سرور روی پورت 8080 اجرا شد...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("خطا در اجرای سرور: ", err)
	}
}