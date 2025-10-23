package paymentmodel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payment represents a single payment record in the database.
type Payment struct {
	PaymentID       primitive.ObjectID	`bson:"_id" json:"paymentID"`
	OrderID  		primitive.ObjectID 	`bson:"orderID" json:"orderID"`
	// PaymentID 		string             	`bson:"paymentId" json:"paymentId"`
	// GatewayOrderID 	string         		`bson:"orderId" json:"orderId"`
	// Signature 		string             	`bson:"signature" json:"signature"`
	Status    		Status             	`bson:"status" json:"status"`
	UserID   		primitive.ObjectID 	`bson:"userID" json:"userID"`
	Price           Price 				`bson:"price" json:"price"`
	CreatedAt 		time.Time			`bson:"createdAt" json:"createdAt"`
	UpdatedAt 		time.Time	 		`bson:"updatedAt" json:"updatedAt"`
}

type Status string
type Price struct {
	TotalAmount   float64 `json:"toal_amount"` // Use float64 for currency precision
	Currency Currency  `json:"currency"`
}
type Currency string
const (
	INR Currency = "INR"
	USD Currency = "USD"
)
const (
	StatusPending Status   = "pending"
	StatusCompleted Status = "completed"
	StatusFailed Status = "failed"
)