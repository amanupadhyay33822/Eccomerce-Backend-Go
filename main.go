package main

import (
	"eccomerce-golang/internal/config"
	"eccomerce-golang/internal/middleware"
	"eccomerce-golang/internal/routes"
	"eccomerce-golang/internal/utils"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	godotenv.Load()
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Apply rate limiting
	r.Use(middleware.RateLimiter())

	config.ConnectDB()

	api := r.Group("/api")
	{
		routes.UserRoutes(api)
		routes.ProductRoutes(api)
		routes.CartRoutes(api)
		routes.OrderRoutes(api)
		routes.PaymentRoutes(api)
	}
	r.GET("/", func(c *gin.Context) {
		utils.JSONResponse(c, 200, true, "Ecommerce API Running ", nil)
	})
	PORT := os.Getenv("PORT")
	fmt.Println(PORT + " is running")
	r.Run(":" + PORT)
}
