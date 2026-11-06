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

// ListUsers mengambil daftar pengguna dengan filter role dan pencarian nama/email.
// Endpoint: GET /api/users (Admin)
func ListUsers(c *gin.Context) {
	var users []models.User
	query := config.DB.Model(&models.User{})

	switch c.Query("role") {
	case "admin":
		query = query.Where("role = ?", models.RoleAdmin)
	case "customer":
		query = query.Where("role = ?", models.RoleCustomer)
	}

	if q := strings.TrimSpace(c.Query("q")); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", like, like)
	}

	if err := query.Order("created_at ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pengguna"})
		return
	}
	if users == nil {
		users = []models.User{}
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUser mengambil data pengguna berdasarkan ID.
// Endpoint: GET /api/users/:id (Admin atau Self)
func GetUser(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	svc := services.NewUserService(config.DB)
	u, err := svc.GetByID(id)
	if err != nil {
		c.JSON(statusForUserError(err), gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

// UpdateUser memperbarui data pengguna.
// Endpoint: PUT /api/users/:id (Admin atau Self)
func UpdateUser(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	actorID := c.GetUint("user_id")
	actorRole, _ := c.Get("role")
	roleStr, _ := actorRole.(string)

	svc := services.NewUserService(config.DB)
	u, err := svc.UpdateUser(actorID, roleStr, id, req)
	if err != nil {
		c.JSON(statusForUserError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pengguna berhasil diperbarui",
		"data":    u,
	})
}

// ChangePassword mengganti kata sandi pengguna.
// Endpoint: PUT /api/users/:id/password (Self saja)
func ChangePassword(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	actorID := c.GetUint("user_id")

	svc := services.NewUserService(config.DB)
	if err := svc.ChangePassword(actorID, id, req); err != nil {
		c.JSON(statusForUserError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kata sandi berhasil diubah"})
}

// DeleteUser menghapus pengguna berdasarkan ID.
// Endpoint: DELETE /api/users/:id (Admin saja)
func DeleteUser(c *gin.Context) {
	id, err := utils.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	actorID := c.GetUint("user_id")

	svc := services.NewUserService(config.DB)
	if err := svc.DeleteUser(actorID, id); err != nil {
		c.JSON(statusForUserError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}

// statusForUserError memetakan error dari service ke HTTP Status Code yang sesuai.
func statusForUserError(err error) int {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, services.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
