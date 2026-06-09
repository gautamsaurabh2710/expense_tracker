package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	if uri == "" {
		log.Fatal("MONGO_URI is required")
	}
	if dbName == "" {
		log.Fatal("DB_NAME is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed:", err)
	}

	fmt.Println("Connected to MongoDB Atlas")
	DB = client.Database(dbName)
}
