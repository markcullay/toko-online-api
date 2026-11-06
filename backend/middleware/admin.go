package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAdmin memastikan pengguna yang mengakses endpoint memiliki Role "admin".
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: khusus admin"})
			return
		}
		c.Next()
	}
}
