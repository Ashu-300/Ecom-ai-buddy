package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() {
	mongoURI := os.Getenv("MONGO_URI")

	ctx  , cancle := context.WithCancel(context.Background())
	defer cancle()

	client , err := mongo.Connect(ctx , options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	log.Printf("✅ Cart Service Connected to MongoDB") ;

	cartCollection = client.Database("supernovaCartDB").Collection("carts")
}