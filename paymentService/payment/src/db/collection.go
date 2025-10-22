package db

import "go.mongodb.org/mongo-driver/mongo"

var paymentCollection *mongo.Collection

func GetPaymentCollection() *mongo.Collection{
	return paymentCollection ;
}
