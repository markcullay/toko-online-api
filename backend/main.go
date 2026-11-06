package main

import (
	"log"
	"os"

	"backend/config"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET wajib di-set di .env")
	}

	config.ConnectDB()

	// Jalankan migrasi otomatis untuk semua model
	if err := config.DB.AutoMigrate(
		&models.Category{},
		&models.User{},
		&models.Product{},
		&models.Sale{},
		&models.SaleItem{},
	); err != nil {
		log.Fatal("Migrasi gagal:", err)
	}

	seedAdmin()

	r := gin.Default()
	r.Use(middleware.SetupCORS())
	routes.Setup(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server berjalan di: http://localhost:" + port)
	r.Run(":" + port)
}

// seedAdmin membuat akun admin default jika belum ada di database.
// Kredensial diambil dari file .env (SEED_ADMIN_EMAIL & SEED_ADMIN_PASSWORD).
func seedAdmin() {
	email := os.Getenv("SEED_ADMIN_EMAIL")
	password := os.Getenv("SEED_ADMIN_PASSWORD")
	if email == "" || password == "" {
		log.Println("SEED_ADMIN_EMAIL/PASSWORD kosong, lewati seeding admin")
		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return
	}

	hash, _ := utils.HashPassword(password)
	admin := models.User{
		Name:     "Admin Toko",
		Email:    email,
		Password: hash,
		Role:     models.RoleAdmin,
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		log.Println("Gagal seed admin:", err)
	} else {
		log.Println("Admin awal dibuat:", email)
	}
}
