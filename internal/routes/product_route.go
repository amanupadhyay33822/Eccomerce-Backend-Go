package routes

import (
	"eccomerce-golang/internal/handlers"
	"eccomerce-golang/internal/middleware"
	"eccomerce-golang/internal/validation"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/products")
	{
		products.GET("", handlers.GetProducts)
		products.GET("/:id", handlers.GetProduct)
		products.PATCH("/:id", handlers.UpdateProduct)
	}

	protected := rg.Group("/products")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/create", middleware.Validate(&validation.ProductInput{}), handlers.CreateProduct)
		protected.DELETE("/:id", handlers.DeleteProduct)
	}
}
