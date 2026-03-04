package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"Name" bson:"name"`
	Description string    `json:"Description" bson:"description"`
	Price       float64   `json:"Price" bson:"price"`
	Stock       int       `json:"Stock" bson:"stock"`
	Category    string    `json:"Category" bson:"category"`
	CreatedAt   time.Time `json:"CreatedAt" bson:"created_at"`
	UpdatedAt   time.Time `json:"UpdatedAt" bson:"updated_at"`
}
