package routes

import (
	"crs-backend/internal/handlers"
	"crs-backend/middleware"
	
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// مسیر احراز هویت کاربران
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register)
		authRoutes.POST("/login", handlers.Login)
	}

	// مسیرهای مرتبط با رزرو بلیط
	bookingRoutes := r.Group("/booking")
	{
		bookingRoutes.GET("/tickets", handlers.GetAvailableTickets)
		bookingRoutes.POST("/reserve", middleware.AuthMiddleware(), handlers.CreateBooking)
		bookingRoutes.GET("/my-bookings", middleware.AuthMiddleware(), handlers.GetUserBookings)
	}

	// مسیرهای مربوط به پروفایل کاربر
	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/profile", handlers.GetUserProfile)
	}

	// مسیرهای مربوط به پرداخت
	paymentRoutes := r.Group("/payment")
	{
		paymentRoutes.POST("/create", handlers.CreatePayment)
		paymentRoutes.GET("/verify", handlers.VerifyPayment)
	}

	// مسیرهای مربوط به مدیریت برای ادمین
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// مدیریت کاربران
		adminRoutes.GET("/users", handlers.GetUsers)
		adminRoutes.GET("/user/:id", handlers.GetUser)
		adminRoutes.DELETE("/user/:id", handlers.DeleteUser)

		// مدیریت بلیط‌ها
		adminRoutes.POST("/add-ticket", handlers.AddTicket)
		adminRoutes.DELETE("/delete-ticket/:id", handlers.DeleteTicket)
		adminRoutes.GET("/all-bookings", handlers.GetAllBookings)
	}

	

	return r

	
}
