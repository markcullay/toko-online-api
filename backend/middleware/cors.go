package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS mengonfigurasi aturan Cross-Origin Resource Sharing
func SetupCORS() gin.HandlerFunc {
	// Ambil dari .env, pisahkan dengan koma jika lebih dari satu
	originsStr := os.Getenv("CORS_ORIGINS")
	if originsStr == "" {
		originsStr = "http://localhost:5500" // fallback default
	}
	origins := strings.Split(originsStr, ",")

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
