package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validate(obj any) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := c.ShouldBindJSON(obj); err != nil {

			errors := make(map[string]string)

			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range ve {

					field := strings.ToLower(fe.Field())

					switch fe.Tag() {

					case "required":
						errors[field] = fe.Field() + " is required"

					case "email":
						errors[field] = "Please provide a valid email address"

					case "min":
						errors[field] = fe.Field() + " must be at least " + fe.Param() + " characters long"

					case "max":
						errors[field] = fe.Field() + " must not exceed " + fe.Param() + " characters"

					case "oneof":
						errors[field] = fe.Field() + " must be one of: " + fe.Param()

					default:
						errors[field] = fe.Field() + " is invalid"
					}
				}
			} else {
				errors["error"] = "Invalid request body"
			}

			c.JSON(400, gin.H{
				"success": false,
				"message": "Validation failed",
				"errors":  errors,
			})

			c.Abort()
			return
		}

		c.Set("validatedBody", obj)
		c.Next()
	}
}
