package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func InitDB()  {
	// MongoDB URI (replace with your own if needed)
	uri := os.Getenv("MONGO_URI")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Error connecting to MongoDB:", err)
	}

	// Ping DB
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Could not ping MongoDB:", err)
	}

	log.Printf("✅ Auth Serive Connected to MongoDB") ;

	UserCollection = client.Database("supernovaAuthDB").Collection("users")

	
}

func CreateUserIndexes(userCollection *mongo.Collection) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    indexModel := mongo.IndexModel{
        Keys:    bson.M{"email": 1},               // index on email
        Options: options.Index().SetUnique(true),  // enforce uniqueness
    }

    _, err := userCollection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        log.Fatal("Failed to create unique index on email:", err)
    } else {
        log.Println("Unique index on email created successfully")
    }
}