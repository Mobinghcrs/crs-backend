package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetDashboardStats برگرداندن اطلاعات آماری داشبورد
func GetDashboardStats(c *gin.Context) {
	stats := map[string]interface{}{
		"totalUsers":   120,
		"totalTickets": 340,
		"totalSales":   5400000,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard Stats",
		"data":    stats,
	})
}
