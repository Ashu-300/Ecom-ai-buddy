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
	mongoUri := os.Getenv("MONGO_URI")
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()

	client , err := mongo.Connect(ctx , options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("error , cannot connect to mongodb" , err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Could not ping MongoDB:", err)
	}

	log.Printf("✅ Payment Service Connected to MongoDB") ;

	paymentCollection = client.Database("SupernovaPaymentDB").Collection("payment")

}