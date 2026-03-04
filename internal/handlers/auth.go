package handlers

import (
	"context"
	"eccomerce-golang/internal/config"
	"eccomerce-golang/internal/models"
	"eccomerce-golang/internal/utils"
	"eccomerce-golang/internal/validation"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(c *gin.Context) {

	body, _ := c.Get("validatedBody")
	input := body.(*validation.RegisterInput)

	collection := config.GetCollection("users")

	// Check duplicate email
	var existing models.User
	err := collection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&existing)
	if err == nil {
		utils.JSONResponse(c, 400, false, "Email already exists", nil)
		return
	}

	// Hash password
	hashedPassword, _ := utils.HashPassword(input.Password)

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      input.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to register user", nil)
		return
	}

	utils.JSONResponse(c, 201, true, "User registered successfully", nil)
}

func Login(c *gin.Context) {

	// Retrieve validated input from middleware
	body, exists := c.Get("validatedBody")
	if !exists {
		utils.JSONResponse(c, 400, false, "Invalid input", nil)
		return
	}
	input := body.(*validation.LoginInput)

	collection := config.GetCollection("users")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)

	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid email or password", nil)
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		utils.JSONResponse(c, 400, false, "Invalid email or password", nil)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)

	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to generate token", nil)
		return
	}

	// ---------------------------
	// Respond
	utils.JSONResponse(c, 200, true, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})

}
