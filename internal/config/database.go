package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DBName string

// ConnectDB initializes MongoDB connection
func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	DBName = os.Getenv("DB_NAME")

	if uri == "" {
		log.Fatal("MONGO_URI not set in environment")
	}

	if DBName == "" {
		log.Fatal("DB_NAME not set in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo connection error:", err)
	}

	// Ping database to verify connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	log.Println("✅ Connected to MongoDB")
}

// GetCollection returns a Mongo collection
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database(DBName).Collection(collectionName)
}
