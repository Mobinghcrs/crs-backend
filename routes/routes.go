package routes

import (
	"crs-backend/internal/handlers"
	"crs-backend/middlewares"
	"crs-backend/internal/repositories"
	"fmt"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, eventRepo repositories.IEventRepository) *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		middlewares.CORS(),
	)

	// � حذف هندلرهای غیرضروری
	eventHandler := handlers.NewEventHandler(eventRepo)
	userHandler := handlers.NewUserHandler(db)

	api := router.Group("/api/v1")
	{
		// � گروه بندی احراز هویت
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", handlers.Register)
			authGroup.POST("/login", handlers.Login)
		}

		// � گروه بندی کاربران با اصلاح نام r به api
		userGroup := api.Group("/users")
		{
			userGroup.GET("", userHandler.GetUsers) // ✅ تغییر به هندلر جدید
			userGroup.GET("/:id", userHandler.GetUser)
			userGroup.DELETE("/:id", userHandler.DeleteUser)
		}

		// � گروه بندی رویدادها
		eventsGroup := api.Group("/events")
		{
			eventsGroup.GET("", eventHandler.GetAllEvents)
			eventsGroup.GET("/:id", eventHandler.GetEventByID)
			
			// � اعتبارسنجی JWT برای عملیات حساس
			authorized := eventsGroup.Use(middlewares.JWTAuth())
			{
				authorized.POST("", eventHandler.CreateEvent)
				authorized.PUT("/:id", eventHandler.UpdateEvent)
				authorized.DELETE("/:id", eventHandler.DeleteEvent)
			}
		}
	}

	// � اندپوینت سلامت سامانه
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"db":     fmt.Sprintf("connected: %v", db != nil),
		})
	})

	printRoutes(router)

	return router
}

func printRoutes(router *gin.Engine) {
	routes := router.Routes()
	fmt.Println("\n=== Registered Routes ===")
	for _, route := range routes {
		fmt.Printf("[%-6s] %s\n", route.Method, route.Path)
	}
	fmt.Println("========================")
}
