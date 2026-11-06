package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SelfOrAdmin mengizinkan akses jika pengguna adalah admin, ATAU jika
// ID pengguna di token sama dengan ID parameter di URL (:id).
func SelfOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role == "admin" {
			c.Next()
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
			return
		}

		paramID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil || uint(paramID) != userID.(uint) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Anda hanya boleh mengakses data Anda sendiri",
			})
			return
		}
		c.Next()
	}
}
