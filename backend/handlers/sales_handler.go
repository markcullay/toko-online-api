package handlers

import (
	"net/http"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// ListSales mengambil riwayat transaksi penjualan terbaru.
// Endpoint: GET /api/sales (Admin)
func ListSales(c *gin.Context) {
	var sales []models.Sale

	err := config.DB.
		Preload("Items").
		Order("created_at DESC").
		Limit(500).
		Find(&sales).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil transaksi"})
		return
	}
	if sales == nil {
		sales = []models.Sale{}
	}
	c.JSON(http.StatusOK, gin.H{"data": sales})
}

// SalesSummary mengambil ringkasan agregasi produk yang terjual beserta total pendapatannya.
// Endpoint: GET /api/sales/summary (Admin)
func SalesSummary(c *gin.Context) {
	type Row struct {
		ProductName string `json:"product"`
		Qty         int64  `json:"qty"`
		Revenue     int64  `json:"revenue"`
	}
	var rows []Row

	err := config.DB.
		Model(&models.SaleItem{}).
		Select("product_name AS product_name, SUM(qty) AS qty, SUM(subtotal) AS revenue").
		Group("product_name").
		Order("qty DESC").
		Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil ringkasan"})
		return
	}
	if rows == nil {
		rows = []Row{}
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
