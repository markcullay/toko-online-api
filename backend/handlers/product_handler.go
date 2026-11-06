package handlers

import (
	"errors"
	"net/http"
	"strings"

	"backend/config"
	"backend/dto"
	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// ListProducts mengambil daftar produk dengan dukungan Pagination, Search, dan Filter Kategori.
// Endpoint: GET /api/products (Publik)
func ListProducts(c *gin.Context) {
	var query dto.ListProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter query tidak valid: " + err.Error()})
		return
	}

	svc := services.NewProductService(config.DB)
	products, meta, err := svc.List(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"meta": meta,
	})
}

// GetProduct mengambil detail satu produk berdasarkan ID.
// Endpoint: GET /api/products/:id (Publik)
func GetProduct(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	svc := services.NewProductService(config.DB)
	p, err := svc.GetByID(id)
	if err != nil {
		c.JSON(statusForProductError(err), gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// CreateProduct menambahkan produk baru ke database.
// Endpoint: POST /api/products (Admin)
func CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	svc := services.NewProductService(config.DB)
	product, err := svc.Create(req)
	if err != nil {
		c.JSON(statusForProductError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil ditambahkan",
		"data":    product,
	})
}

// UpdateProduct memperbarui detail produk berdasarkan ID.
// Endpoint: PUT /api/products/:id (Admin)
func UpdateProduct(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	svc := services.NewProductService(config.DB)
	product, err := svc.Update(id, req)
	if err != nil {
		c.JSON(statusForProductError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diperbarui",
		"data":    product,
	})
}

// DeleteProduct menghapus produk berdasarkan ID.
// Endpoint: DELETE /api/products/:id (Admin)
func DeleteProduct(c *gin.Context) {
	result := config.DB.Delete(&models.Product{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}

// statusForProductError memetakan error dari service ke HTTP Status Code yang sesuai.
func statusForProductError(err error) int {
	switch {
	case errors.Is(err, services.ErrProductNotFound):
		return http.StatusNotFound
	case errors.Is(err, services.ErrProductConflict):
		return http.StatusConflict
	default:
		if strings.Contains(err.Error(), "kategori tidak valid") {
			return http.StatusBadRequest
		}
		return http.StatusBadRequest
	}
}
