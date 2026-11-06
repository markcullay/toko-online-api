package middleware

import (
	"net/http"
	"strings"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

// Auth memvalidasi JWT Token dari header Authorization.
// Jika valid, data user (ID dan Role) disimpan ke dalam Context Gin.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluarsa"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
