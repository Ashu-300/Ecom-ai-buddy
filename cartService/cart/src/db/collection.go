package db

import "go.mongodb.org/mongo-driver/mongo"

var cartCollection *mongo.Collection

func GetCartCollection() *mongo.Collection {
	return cartCollection
}