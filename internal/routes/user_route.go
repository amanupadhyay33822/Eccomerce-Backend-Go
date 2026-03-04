package routes

import (
	"eccomerce-golang/internal/handlers"
	"eccomerce-golang/internal/middleware"
	"eccomerce-golang/internal/validation"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {

	auth := rg.Group("/auth")
	{
		auth.POST("/register",
			middleware.Validate(&validation.RegisterInput{}),
			handlers.Register,
		)

		auth.POST("/login",
			middleware.Validate(&validation.LoginInput{}),
			handlers.Login,
		)

	}

}
