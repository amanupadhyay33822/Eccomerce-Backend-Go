package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID        primitive.ObjectID `bson:"_id"`
	OrderID   primitive.ObjectID `bson:"order_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Amount    float64            `bson:"amount"`
	Currency  string             `bson:"currency"`
	Status    string             `bson:"status"`     // pending, success, failed
	GatewayID string             `bson:"gateway_id"` // razorpay/stripe payment id
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentGateway interface {
	CreatePayment(amount float64, currency string) (string, error)
	VerifyPayment(data map[string]string) error
}