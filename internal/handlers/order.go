package handlers

import (
	"eccomerce-golang/internal/config"
	"eccomerce-golang/internal/models"
	"eccomerce-golang/internal/utils"
	"eccomerce-golang/internal/validation"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateOrder - place a new order
func CreateOrder(c *gin.Context) {
	body, exist := c.Get("validatedBody")
	if !exist {
		utils.JSONResponse(c, 400, false, "Invalid input", nil)
		return
	}
	input := body.(*validation.CreateOrderInput)
	userID := utils.GetUserID(c)
	collection := config.GetCollection("orders")

	// Calculate total price
	total := 0.0
	var items []models.OrderItems
	productCollection := config.GetCollection("products")

	for _, i := range input.Items {
		productID, err := primitive.ObjectIDFromHex(i.ProductID)
		if err != nil {
			utils.JSONResponse(c, 400, false, "Invalid product ID", nil)
			return
		}

		// Fetch product price from database
		var product models.Product
		err = productCollection.FindOne(c.Request.Context(), bson.M{"_id": productID}).Decode(&product)
		if err != nil {
			utils.JSONResponse(c, 404, false, "Product not found", nil)
			return
		}

		// Check stock availability
		if product.Stock < i.Quantity {
			utils.JSONResponse(c, 400, false, "Insufficient stock for product", nil)
			return
		}

		// Decrease stock
		_, err = productCollection.UpdateOne(
			c.Request.Context(),
			bson.M{"_id": productID},
			bson.M{"$inc": bson.M{"stock": -i.Quantity}},
		)
		if err != nil {
			utils.JSONResponse(c, 500, false, "Failed to update stock", nil)
			return
		}

		items = append(items, models.OrderItems{
			ProductID: productID,
			Quantity:  i.Quantity,
			Price:     product.Price,
		})
		total += float64(i.Quantity) * product.Price
	}

	order := models.Order{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		Items:      items,
		TotalPrice: total,
		Status:     "pending",
		CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err := collection.InsertOne(c.Request.Context(), order)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to create order", nil)
		return
	}

	utils.JSONResponse(c, 201, true, "Order placed successfully", order)
}

// GetOrders - get all orders of logged-in user
func GetOrders(c *gin.Context) {
	userID := utils.GetUserID(c)
	collection := config.GetCollection("orders")

	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(c.Request.Context(), filter)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to fetch orders", nil)
		return
	}
	defer cursor.Close(c.Request.Context())

	var orders []models.Order
	if err := cursor.All(c.Request.Context(), &orders); err != nil {
		utils.JSONResponse(c, 500, false, "Failed to decode orders", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Orders fetched successfully", orders)
}

// GetOrder - get single order by ID
func GetOrder(c *gin.Context) {
	objectID, err := utils.GetParamObjectID(c, "id")
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid order ID", nil)
		return
	}

	userID := utils.GetUserID(c)
	collection := config.GetCollection("orders")

	var order models.Order
	err = collection.FindOne(c.Request.Context(), bson.M{"_id": objectID, "user_id": userID}).Decode(&order)
	if err != nil {
		utils.JSONResponse(c, 404, false, "Order not found", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Order fetched successfully", order)
}

// UpdateOrderStatus - update order status (admin or user)
func UpdateOrderStatus(c *gin.Context) {
	objectID, err := utils.GetParamObjectID(c, "id")
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid order ID", nil)
		return
	}

	body, exist := c.Get("validatedBody")
	if !exist {
		utils.JSONResponse(c, 400, false, "Invalid input", nil)
		return
	}
	input := body.(*validation.UpdateOrderStatusInput)
	collection := config.GetCollection("orders")

	update := bson.M{
		"$set": bson.M{
			"status":     input.Status,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(c.Request.Context(), bson.M{"_id": objectID}, update)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to update order", nil)
		return
	}

	if result.MatchedCount == 0 {
		utils.JSONResponse(c, 404, false, "Order not found", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Order status updated successfully", nil)
}

// DeleteOrder - cancel an order
func DeleteOrder(c *gin.Context) {
	objectID, err := utils.GetParamObjectID(c, "id")
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid order ID", nil)
		return
	}

	collection := config.GetCollection("orders")

	result, err := collection.DeleteOne(c.Request.Context(), bson.M{"_id": objectID})
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to delete order", nil)
		return
	}

	if result.DeletedCount == 0 {
		utils.JSONResponse(c, 404, false, "Order not found", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Order cancelled successfully", nil)
}
