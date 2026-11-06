package handlers

import (
	"net/http"
	"strings"

	"backend/config"
	"backend/dto"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// Register mendaftarkan pengguna baru sebagai customer.
// Endpoint: POST /api/auth/register (Publik)
func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Cek email sudah terdaftar?
	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	// Register publik SELALU jadi customer.
	// Admin pertama dibuat manual / via seeder.
	user := models.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    req.Email,
		Password: hash,
		Role:     models.RoleCustomer,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Akun berhasil dibuat",
		"user":    user,
	})
}

// Login memverifikasi kredensial pengguna dan mengembalikan JWT Token.
// Endpoint: POST /api/auth/login (Publik)
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Pesan sengaja generik — jangan bocorkan mana yang salah
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Me mengambil data profil pengguna yang sedang login.
// Endpoint: GET /api/auth/me (Butuh Login)
func Me(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
