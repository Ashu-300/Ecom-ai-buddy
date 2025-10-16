package db

import "go.mongodb.org/mongo-driver/mongo"

var orderCollection *mongo.Collection

func GetOrderCollection() *mongo.Collection {
	return orderCollection
}