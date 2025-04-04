package routes

import (
	"crs-backend/cp-backend/handlers"
	"crs-backend/cp-backend/middleware"
	"github.com/gin-gonic/gin"
)

// SetupCPRoutes تنظیم مسیرهای کنترل پنل (CP) را بر عهده دارد
func SetupCPRoutes(r *gin.Engine) {
	cp := r.Group("/cp")
	cp.Use(middleware.AdminAuthMiddleware())

	// مسیر داشبورد
	cp.GET("/dashboard", handlers.GetDashboardStats)

	// مسیرهای مدیریت کاربران
	users := cp.Group("/users")
	{
		users.GET("/", handlers.ListUsers)
		users.POST("/", handlers.CreateUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}

	// مسیرهای مدیریت بلیط‌ها (مثلاً اگر بخواهید بلیط‌ها را مدیریت کنید)
	tickets := cp.Group("/tickets")
	{
		tickets.GET("/", handlers.ListTickets)
		tickets.POST("/", handlers.CreateTicket)
		tickets.PUT("/:id", handlers.UpdateTicket)
		tickets.DELETE("/:id", handlers.DeleteTicket)
	}

	// مسیرهای مدیریت پروازها
	flights := cp.Group("/flights")
	{
		flights.GET("/", handlers.ListFlights)       // دریافت لیست پروازها (اختیاری)
		flights.POST("/", handlers.CreateFlight)       // ایجاد/افزودن پرواز
		flights.PUT("/:id", handlers.UpdateFlight)       // ویرایش پرواز
		flights.DELETE("/:id", handlers.DeleteFlight)    // حذف پرواز
	}
}
