package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() {
	uri := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client , err := mongo.Connect(ctx , options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("✅ sellerDashboard Service Connected to MongoDB!")

	sellerUserCollection = client.Database("SupernovaSellerDashboardDB").Collection("user")
	sellerOrderCollection = client.Database("SupernovaSellerDashboardDB").Collection("order")
	sellerPaymentCollection = client.Database("SupernovaSellerDashboardDB").Collection("payment")
	sellerProductCollection = client.Database("SupernovaSellerDashboardDB").Collection("product")
}