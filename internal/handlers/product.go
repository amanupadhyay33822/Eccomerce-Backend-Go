package handlers

import (
	"context"
	"eccomerce-golang/internal/config"
	"eccomerce-golang/internal/models"
	"eccomerce-golang/internal/utils"
	"eccomerce-golang/internal/validation"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(c *gin.Context) {
	body, _ := c.Get("validatedBody")
	input := body.(*validation.ProductInput)

	collection := config.GetCollection("products")

	newProduct := models.Product{
		ID:          primitive.NewObjectID(),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		Category:    input.Category,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
fmt.Println(`product` ,  input)
	_, err := collection.InsertOne(context.Background(), newProduct)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to create product", nil)
		return
	}

	utils.JSONResponse(c, 201, true, "Product created successfully", newProduct)
}

func GetProducts(c *gin.Context) {
	fmt.Println(`bdbnsmnd`)
	pagination := utils.GetPaginationParams(c)
	filter := utils.GetProductFilter(c)

	collection := config.GetCollection("products")

	opts := utils.GetFindOptions(pagination)
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to fetch products", nil)
		return
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	if err = cursor.All(context.Background(), &products); err != nil {
		utils.JSONResponse(c, 500, false, "Failed to decode products", nil)
		return
	}

	total, _ := collection.CountDocuments(context.Background(), filter)

	utils.JSONResponse(c, 200, true, "Productsss fetched successfully", gin.H{
		"products":   products,
		"pagination": gin.H{
			"page":       pagination.Page,
			"limit":      pagination.Limit,
			"total":      total,
			"totalPages": (total + int64(pagination.Limit) - 1) / int64(pagination.Limit),
		},
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid product ID", nil)
		return
	}

	collection := config.GetCollection("products")

	var product models.Product
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		utils.JSONResponse(c, 404, false, "Product not found", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Product fetched successfully", product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid product ID", nil)
		return
	}

	var input validation.UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONResponse(c, 400, false, "Invalid input", nil)
		return
	}

	collection := config.GetCollection("products")

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	setFields := update["$set"].(bson.M)

	if input.Name != nil {
		setFields["name"] = *input.Name
	}
	if input.Description != nil {
		setFields["description"] = *input.Description
	}
	if input.Price != nil {
		setFields["price"] = *input.Price
	}
	if input.Stock != nil {
		setFields["stock"] = *input.Stock
	}
	if input.Category != nil {
		setFields["category"] = *input.Category
	}

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to update product", nil)
		return
	}

	if result.MatchedCount == 0 {
		utils.JSONResponse(c, 404, false, "Product not found", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Product updated successfully", nil)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid product ID", nil)
		return
	}

	collection := config.GetCollection("products")

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to delete product", nil)
		return
	}

	if result.DeletedCount == 0 {
		utils.JSONResponse(c, 404, false, "Product not found", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Product deleted successfully", nil)
}
