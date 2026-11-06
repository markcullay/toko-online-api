package handlers

import (
	"errors"
	"net/http"

	"backend/config"
	"backend/dto"
	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// ListCategories mengambil semua daftar kategori.
// Endpoint: GET /api/categories (Publik)
func ListCategories(c *gin.Context) {
	svc := services.NewCategoryService(config.DB)
	categories, err := svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateCategory menambahkan kategori baru ke database.
// Endpoint: POST /api/categories (Admin)
func CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	svc := services.NewCategoryService(config.DB)
	category, err := svc.Create(req)
	if err != nil {
		c.JSON(statusForCategoryError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil ditambahkan",
		"data":    category,
	})
}

// UpdateCategory memperbarui data kategori berdasarkan ID.
// Endpoint: PUT /api/categories/:id (Admin)
func UpdateCategory(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	svc := services.NewCategoryService(config.DB)
	category, err := svc.Update(id, req)
	if err != nil {
		c.JSON(statusForCategoryError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil diperbarui",
		"data":    category,
	})
}

// DeleteCategory menghapus kategori berdasarkan ID.
// Akan gagal jika masih ada produk yang menggunakan kategori ini.
// Endpoint: DELETE /api/categories/:id (Admin)
func DeleteCategory(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	svc := services.NewCategoryService(config.DB)
	if err := svc.Delete(id); err != nil {
		c.JSON(statusForCategoryError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}

// statusForCategoryError memetakan error dari service ke HTTP Status Code yang sesuai.
func statusForCategoryError(err error) int {
	switch {
	case errors.Is(err, services.ErrCategoryNotFound):
		return http.StatusNotFound
	case errors.Is(err, services.ErrCategoryConflict):
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}
