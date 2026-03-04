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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func calculateCartTotal(c *gin.Context, userID primitive.ObjectID) {
	cartCollection := config.GetCollection("carts")
	var cart models.Cart
	err := cartCollection.FindOne(c.Request.Context(), bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return
	}

	productCollection := config.GetCollection("products")
	var totalPrice float64
	for _, item := range cart.Items {
		var product models.Product
		err := productCollection.FindOne(c.Request.Context(), bson.M{"_id": item.ProductID}).Decode(&product)
		if err == nil {
			totalPrice += product.Price * float64(item.Quantity)
		}
	}

	update := bson.M{"$set": bson.M{"total_price": totalPrice}}
	cartCollection.UpdateOne(c.Request.Context(), bson.M{"user_id": userID}, update)
}

func AddToCart(c *gin.Context) {
	body, exist := c.Get("validatedBody")
	if !exist {
		utils.JSONResponse(c, 400, false, "Invalid input", nil)
		return
	}

	input := body.(*validation.AddToCartInput)
	userIDStr := c.MustGet("user_id").(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid user ID", nil)
		return
	}

	filter := bson.M{"user_id": userID, "items.product_id": input.ProductID}

	//check if the product is already in the cart
	var existingCartItem models.Cart
	collections := config.GetCollection("carts")
	err = collections.FindOne(c.Request.Context(), filter).Decode(&existingCartItem)
	if err == nil {
		// If the product is already in the cart, update the quantity
		update := bson.M{
			"$inc": bson.M{"items.$.quantity": input.Quantity},
			"$set": bson.M{"updated_at": time.Now()},
		}
		_, err := collections.UpdateOne(c.Request.Context(), filter, update)
		if err != nil {
			utils.JSONResponse(c, 500, false, "Failed to update cart", nil)
			return
		}
		calculateCartTotal(c, userID)
		utils.JSONResponse(c, 200, true, "Cart updated successfully", nil)
		return
	}

	// If the product is not in the cart, add it as a new item
	productID, err := primitive.ObjectIDFromHex(input.ProductID)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid product ID", nil)
		return
	}

	newCartItem := models.CartItem{
		ProductID: productID,
		Quantity:  input.Quantity,
	}
	update := bson.M{
		"$push": bson.M{"items": newCartItem},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	opts := options.Update().SetUpsert(true)
	_, err = collections.UpdateOne(c.Request.Context(), bson.M{"user_id": userID}, update, opts)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to add to cart", nil)
		return
	}
	calculateCartTotal(c, userID)
	utils.JSONResponse(c, 200, true, "Item added to cart successfully", nil)

}

func GetCart(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid user ID", nil)
		return
	}

	collection := config.GetCollection("carts")

	var cart models.Cart
	err = collection.FindOne(c.Request.Context(), bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		utils.JSONResponse(c, 404, false, "Cart not found", nil)
		return
	}

	utils.JSONResponse(c, 200, true, "Cart retrieved successfully", cart)
}

func UpdateCartItem(c *gin.Context) {
	body, exist := c.Get("validatedBody")
	if !exist {
		utils.JSONResponse(c, 400, false, "Invalid input", nil)
		return
	}
	input := body.(*validation.UpdateCartItemInput)

	userIDStr := c.MustGet("user_id").(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid user ID", nil)
		return
	}

	filter := bson.M{"user_id": userID, "items.product_id": input.ProductID}

	update := bson.M{
		"$set": bson.M{
			"items.$.quantity": input.Quantity,
			"updated_at":       time.Now(),
		},
	}

	collection := config.GetCollection("carts")
	result, err := collection.UpdateOne(c.Request.Context(), filter, update)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to update cart item", nil)
		return
	}
	if result.MatchedCount == 0 {
		utils.JSONResponse(c, 404, false, "Cart item not found", nil)
		return
	}
	calculateCartTotal(c, userID)
	utils.JSONResponse(c, 200, true, "Cart item updated successfully", nil)
}

func RemoveCartItem(c *gin.Context) {
	productID := c.Param("product_id")
	userIDStr := c.MustGet("user_id").(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid user ID", nil)
		return
	}

	productObjID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		utils.JSONResponse(c, 400, false, "Invalid product ID", nil)
		return
	}

	collection := config.GetCollection("carts")
	update := bson.M{
		"$pull": bson.M{"items": bson.M{"product_id": productObjID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err = collection.UpdateOne(c.Request.Context(), bson.M{"user_id": userID}, update)
	if err != nil {
		utils.JSONResponse(c, 500, false, "Failed to remove item from cart", nil)
		return
	}
	calculateCartTotal(c, userID)
	utils.JSONResponse(c, 200, true, "Item removed from cart successfully", nil)
}
