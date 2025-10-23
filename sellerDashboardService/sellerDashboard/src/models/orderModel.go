package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order represents an e-commerce order.
type Order struct {
	// Use primitive.ObjectID for the MongoDB primary key (_id)
	OrderID         primitive.ObjectID `json:"orderId" bson:"_id" binding:"required"`
	UserID          primitive.ObjectID `json:"userId" bson:"userId" binding:"required"`
	Items           []Item    	    	`json:"items" bson:"items" binding:"required"`
	TotalPrice     	Price            	`json:"totalPrice" bson:"totalPrice" binding:"required"`
	Status          OrderStatus        `json:"status" bson:"status"`
	Address			Address    		   `json:"address" bson:"address" binding:"required"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt" binding:"required"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt" binding:"required"`
}

// OrderItem represents a single product within an order.
type Item struct {
	ProductID primitive.ObjectID `bson:"productId" json:"productId"`
	Price     Price            `bson:"price" json:"price"`
	Quantity  int                `bson:"quantity" json:"quantity" binding:"required,min=1"`
}

type Price struct {
	Amount   float64 `json:"amount" bson:"amount" binding:"required"`
	Currency Currency  `json:"currency" bson:"currency" binding:"required"`
}

type Currency string
const (
	USD Currency = "USD"
	INR Currency = "INR"
)

// ShippingAddress represents the delivery location for an order.


// OrderStatus is an enum for the current state of an order.
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)