package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api")

	// ---------- Public ----------
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Products & Categories (Baca)
	api.GET("/products", handlers.ListProducts)
	api.GET("/products/:id", handlers.GetProduct)
	api.GET("/categories", handlers.ListCategories)

	// ---------- Butuh login ----------
	authn := api.Group("")
	authn.Use(middleware.Auth())
	{
		authn.GET("/auth/me", handlers.Me)
		authn.POST("/checkout", handlers.Checkout)
		authn.PUT("/users/:id/password", handlers.ChangePassword)
	}

	// ---------- Self atau admin ----------
	selfOrAdmin := api.Group("")
	selfOrAdmin.Use(middleware.Auth(), middleware.SelfOrAdmin())
	{
		selfOrAdmin.GET("/users/:id", handlers.GetUser)
		selfOrAdmin.PUT("/users/:id", handlers.UpdateUser)
	}

	// ---------- Butuh admin ----------
	admin := api.Group("")
	admin.Use(middleware.Auth(), middleware.RequireAdmin())
	{
		// Category Management
		admin.POST("/categories", handlers.CreateCategory)
		admin.PUT("/categories/:id", handlers.UpdateCategory)
		admin.DELETE("/categories/:id", handlers.DeleteCategory)

		// Product Management
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)

		// Sales & Stats
		admin.GET("/sales", handlers.ListSales)
		admin.GET("/sales/summary", handlers.SalesSummary)
		admin.GET("/stats", handlers.Stats)

		// Users
		admin.GET("/users", handlers.ListUsers)
		admin.DELETE("/users/:id", handlers.DeleteUser)
	}
}
