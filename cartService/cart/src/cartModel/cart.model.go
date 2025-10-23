package cartmodel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ProductID 	primitive.ObjectID `bson:"productId" json:"productId"`
	Price 		Price				`bson:"price" json:"price"`
	Quantity  	int                `bson:"quantity" json:"quantity" binding:"required,min=1"`
}
type Price struct {
	Amount  	float64 `bson:"amount,omitempty" json:"amount,omitempty"`
	Currency	Currency  `bson:"currency,omitempty" json:"currency,omitempty"`
}
type Currency string

const (
	USD Currency = "USD"
	INR Currency = "INR"
)

type Cart struct {
	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
	Items      []Item             `bson:"items" json:"items"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
