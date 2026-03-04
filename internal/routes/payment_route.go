package routes

import (
	"eccomerce-golang/internal/handlers"
	"eccomerce-golang/internal/middleware"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(rg *gin.RouterGroup) {
	payments := rg.Group("/payments")
	{
		// Create payment (requires auth)
		payments.POST("/:order_id", middleware.AuthMiddleware(), handlers.CreatePayment)
		
		// Verify payment (requires auth)
		payments.POST("/verify", middleware.AuthMiddleware(), handlers.VerifyPayment)
		
		// Webhook (no auth - called by Razorpay)
		payments.POST("/webhook", handlers.PaymentWebhook)
	}
}
