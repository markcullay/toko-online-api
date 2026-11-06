package handlers

import (
	"net/http"

	"backend/config"
	"backend/dto"
	"backend/services"

	"github.com/gin-gonic/gin"
)

// Checkout memproses pembelian produk oleh pengguna.
// Endpoint: POST /api/checkout (Login)
func Checkout(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	svc := services.NewCheckoutService(config.DB)
	sale, err := svc.Checkout(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pembayaran berhasil",
		"data":    sale,
	})
}
