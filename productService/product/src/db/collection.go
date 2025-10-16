package db

import "go.mongodb.org/mongo-driver/mongo"

var productCollection *mongo.Collection

func GetProductCollection() *mongo.Collection {
	return productCollection
}
