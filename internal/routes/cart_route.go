package routes

import (
	"eccomerce-golang/internal/handlers"
	"eccomerce-golang/internal/middleware"
	"eccomerce-golang/internal/validation"

	"github.com/gin-gonic/gin"
)

func CartRoutes(rg *gin.RouterGroup) {
	cart := rg.Group("/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.POST("/add", middleware.Validate(&validation.AddToCartInput{}), handlers.AddToCart)
		cart.GET("", handlers.GetCart)
		cart.PUT("/update", middleware.Validate(&validation.UpdateCartItemInput{}), handlers.UpdateCartItem)
		cart.DELETE("/remove/:product_id", handlers.RemoveCartItem)
	}
}
