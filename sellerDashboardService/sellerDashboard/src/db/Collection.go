package db

import "go.mongodb.org/mongo-driver/mongo"

var sellerUserCollection *mongo.Collection
var sellerOrderCollection *mongo.Collection
var sellerPaymentCollection *mongo.Collection
var sellerProductCollection *mongo.Collection

func GetSellerUserCollection() *mongo.Collection {
	return sellerUserCollection
}

func GetSellerOrderCollection()*mongo.Collection{
	return sellerOrderCollection
}

func GetSellerPaymentCollection()*mongo.Collection{
	return sellerPaymentCollection
}

func GetSellerProductCollection() *mongo.Collection{
	return sellerProductCollection 
}