package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItems struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items      []OrderItems       `bson:"items" json:"items"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt  primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
