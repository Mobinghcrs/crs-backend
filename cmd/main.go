package main

import (
	"crs-backend/internal/database"
	"crs-backend/routes"                      // مسیرهای اصلی CRS
	cpRoutes "crs-backend/cp-backend/routes"   // مسیرهای CP با alias برای جلوگیری از تداخل
	//"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	// اضافه کردن مسیرهای CRS
	router := routes.SetupRouter()


	// اضافه کردن مسیرهای CP
	cpRoutes.SetupCPRoutes(router)

	// اجرای سرور روی پورت 8080
	router.Run(":8080")
}
