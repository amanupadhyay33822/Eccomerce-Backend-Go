package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"eccomerce-golang/internal/config"
	"eccomerce-golang/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//////////////////////////////////////////////////////////////////
// 🔐 Helper: Verify Signature
//////////////////////////////////////////////////////////////////

func verifySignature(body []byte, signature string, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

//////////////////////////////////////////////////////////////////
// 💳 1️⃣ Create Payment (Create Razorpay Order)
//////////////////////////////////////////////////////////////////

func CreatePayment(c *gin.Context) {

	orderIDStr := c.Param("order_id")
	orderID, err := primitive.ObjectIDFromHex(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	userID := c.MustGet("user_id").(primitive.ObjectID)

	orderCollection := config.GetCollection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order bson.M
	err = orderCollection.FindOne(ctx, bson.M{
		"_id":     orderID,
		"user_id": userID,
	}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Safe amount extraction
	amountValue, ok := order["total_price"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid order amount"})
		return
	}

	amount, ok := amountValue.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "amount type error"})
		return
	}

	// Razorpay Client
	key := os.Getenv("RAZORPAY_KEY")
	secret := os.Getenv("RAZORPAY_SECRET")

	client := razorpay.NewClient(key, secret)

	data := map[string]interface{}{
		"amount":   int(amount * 100), // convert to paisa
		"currency": "INR",
		"receipt":  orderIDStr,
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gateway error"})
		return
	}

	gatewayID := body["id"].(string)

	paymentCollection := config.GetCollection("payments")

	payment := models.Payment{
		ID:        primitive.NewObjectID(),
		OrderID:   orderID,
		UserID:    userID,
		Amount:    amount,
		Currency:  "INR",
		Status:    "pending",
		GatewayID: gatewayID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = paymentCollection.InsertOne(ctx, payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gateway_order_id": gatewayID,
		"amount":           amount,
		"currency":         "INR",
	})
}

//////////////////////////////////////////////////////////////////
// ✅ 2️⃣ Verify Payment (Frontend Method)
//////////////////////////////////////////////////////////////////

func VerifyPayment(c *gin.Context) {

	var input struct {
		OrderID        string `json:"order_id" binding:"required"`
		PaymentID      string `json:"payment_id" binding:"required"`
		Signature      string `json:"signature" binding:"required"`
		GatewayOrderID string `json:"gateway_order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	secret := os.Getenv("RAZORPAY_SECRET")

	// Verify signature
	message := input.GatewayOrderID + "|" + input.PaymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(input.Signature)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	paymentCollection := config.GetCollection("payments")

	// Idempotency check
	var existing models.Payment
	err := paymentCollection.FindOne(ctx, bson.M{
		"gateway_id": input.GatewayOrderID,
	}).Decode(&existing)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment not found"})
		return
	}

	if existing.Status == "success" {
		c.JSON(http.StatusOK, gin.H{"message": "already verified"})
		return
	}

	// Update payment
	_, err = paymentCollection.UpdateOne(ctx,
		bson.M{"gateway_id": input.GatewayOrderID},
		bson.M{
			"$set": bson.M{
				"status":     "success",
				"payment_id": input.PaymentID,
				"updated_at": time.Now(),
			},
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment"})
		return
	}

	// Update order
	orderID, _ := primitive.ObjectIDFromHex(input.OrderID)
	orderCollection := config.GetCollection("orders")

	_, err = orderCollection.UpdateOne(ctx,
		bson.M{"_id": orderID},
		bson.M{
			"$set": bson.M{
				"status":     "confirmed",
				"updated_at": time.Now(),
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment verified successfully"})
}

//////////////////////////////////////////////////////////////////
// 🌐 3️⃣ Webhook Method (Backend Verification)
//////////////////////////////////////////////////////////////////

func PaymentWebhook(c *gin.Context) {

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	signature := c.GetHeader("X-Razorpay-Signature")
	secret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")

	if !verifySignature(bodyBytes, signature, secret) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	var payload struct {
		Event   string `json:"event"`
		Payload struct {
			Payment struct {
				Entity struct {
					ID      string `json:"id"`
					OrderID string `json:"order_id"`
					Status  string `json:"status"`
				} `json:"entity"`
			} `json:"payment"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if payload.Event != "payment.captured" {
		c.JSON(http.StatusOK, gin.H{"message": "ignored"})
		return
	}

	gatewayOrderID := payload.Payload.Payment.Entity.OrderID
	paymentID := payload.Payload.Payment.Entity.ID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	paymentCollection := config.GetCollection("payments")

	var existing models.Payment
	err = paymentCollection.FindOne(ctx, bson.M{
		"gateway_id": gatewayOrderID,
	}).Decode(&existing)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment not found"})
		return
	}

	if existing.Status == "success" {
		c.JSON(http.StatusOK, gin.H{"message": "already processed"})
		return
	}

	// Update payment
	_, err = paymentCollection.UpdateOne(ctx,
		bson.M{"gateway_id": gatewayOrderID},
		bson.M{
			"$set": bson.M{
				"status":     "success",
				"payment_id": paymentID,
				"updated_at": time.Now(),
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment"})
		return
	}

	// Update order
	orderCollection := config.GetCollection("orders")

	_, err = orderCollection.UpdateOne(ctx,
		bson.M{"_id": existing.OrderID},
		bson.M{
			"$set": bson.M{
				"status":     "confirmed",
				"updated_at": time.Now(),
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment confirmed"})
}
