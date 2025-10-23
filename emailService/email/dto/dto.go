package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type JsonUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PaymentData struct {
	ReceiverMail string 			`json:"receiverMail"`
	PaymentID  	primitive.ObjectID `json:"paymentID"`
	OrderID		primitive.ObjectID `json:"orderID"`
	Amount 		float64				`json:"amount"`
	Currency    string				`json:"currency"`
}

type ProductData struct {
    ReceiverMail string			`json:"receiverMail"`
    ProductName  string			`json:"receiverMail"`
    ProductID    string			`json:"receiverMail"`
    Price        float64		`json:"receiverMail"`
    Currency     string			`json:"receiverMail"`
}


type OrderData struct {
	ReceiverMail string
	OrderID      string
	TotalAmount  float64
	Currency     string
}
