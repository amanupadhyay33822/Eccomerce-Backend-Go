package routes

import (
	"eccomerce-golang/internal/handlers"
	"eccomerce-golang/internal/middleware"
	"eccomerce-golang/internal/validation"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("/create", middleware.Validate(&validation.CreateOrderInput{}), handlers.CreateOrder)
		orders.GET("", handlers.GetOrders)
		orders.GET("/:id", handlers.GetOrder)
		orders.PUT("/:id/status", middleware.Validate(&validation.UpdateOrderStatusInput{}), handlers.UpdateOrderStatus)
		orders.DELETE("/:id", handlers.DeleteOrder)
	}
}
