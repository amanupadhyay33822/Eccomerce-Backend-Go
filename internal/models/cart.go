package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
}

type Cart struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items      []CartItem         `bson:"items" json:"items"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
