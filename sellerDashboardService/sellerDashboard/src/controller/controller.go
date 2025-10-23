package controller

import (
	"context"
	"log"
	"supernova/sellerDashboardService/sellerDashboard/src/db"
	"supernova/sellerDashboardService/sellerDashboard/src/models"
	"time"
)

func CreateUser(user models.User) {
	userCollection := db.GetSellerUserCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := userCollection.InsertOne(ctx , user)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}

func CreateProduct(product models.Product){
	productCollection := db.GetSellerProductCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := productCollection.InsertOne(ctx , product)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}


func CreateOrder(order models.Order){
	orderCollection := db.GetSellerOrderCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := orderCollection.InsertOne(ctx , order)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}

func CreatePayment(payment models.Order){
	paymentCollection := db.GetSellerPaymentCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := paymentCollection.InsertOne(ctx , payment)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}