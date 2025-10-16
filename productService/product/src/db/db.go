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

var DB *mongo.Database

func InItDB() {
	uri := os.Getenv("MONGO_URI")
	
	ctx , cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	
	client , err := mongo.Connect(ctx , options.Client().ApplyURI(uri)) 
	
	if err != nil {
		log.Fatal("❌ Error connecting to MongoDB:", err)
	}

	// Ping DB
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Could not ping MongoDB:", err)
	}

	log.Printf("✅ Product Service Connected to MongoDB") ;

	productCollection = client.Database("SupernovaProductDB").Collection("products")
	
	err = CreateProductIndex(productCollection)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

}


func CreateProductIndex(collection *mongo.Collection) error {
    ctx := context.Background()
    
    indexModel := mongo.IndexModel{
        Keys: bson.D{
            {Key: "title", Value: 1},
            {Key: "description", Value: 1},
        },
        Options: options.Index().SetName("title_description_index"),
    }

    _, err := collection.Indexes().CreateOne(ctx, indexModel)
    return err
}