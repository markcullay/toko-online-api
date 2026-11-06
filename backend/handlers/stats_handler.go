package handlers

import (
	"net/http"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// Stats mengambil data statistik ringkasan untuk dashboard admin.
// Endpoint: GET /api/stats (Admin)
func Stats(c *gin.Context) {
	var totalProducts, totalUsers, totalSold int64
	var totalRevenue int64

	config.DB.Model(&models.Product{}).Count(&totalProducts)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleCustomer).Count(&totalUsers)

	config.DB.Model(&models.SaleItem{}).Select("COALESCE(SUM(qty),0)").Scan(&totalSold)
	config.DB.Model(&models.SaleItem{}).Select("COALESCE(SUM(subtotal),0)").Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"total_products": totalProducts,
		"total_users":    totalUsers,
		"total_sold":     totalSold,
		"total_revenue":  totalRevenue,
	})
}
