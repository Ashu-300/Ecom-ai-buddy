package dto

import (
	ordermodel "supernova/orderService/order/src/orderModel"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Price defines the price structure with amount and currency
// type Price struct {
// 	Amount   float64 `bson:"amount" json:"amount"`
// 	Currency string  `bson:"currency" json:"currency"`
// }

// // Item defines a single product inside the cart
// type Item struct {
// 	ProductID primitive.ObjectID `bson:"productId" json:"productId"`
// 	Price     Price               `bson:"price" json:"price"`
// 	Quantity  int                 `bson:"quantity" json:"quantity"`
// }

// Cart defines the entire cart structure
type Cart struct {
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Items     []ordermodel.Item             `bson:"items" json:"items"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// ResponseCart defines the response structure sent back to the client
type ResponseCart struct {
	Message string `json:"message"`
	Cart    Cart   `json:"cart"`
}



type OrderData struct {
	ReceiverMail string
	OrderID      primitive.ObjectID
	TotalAmount  float64
	Currency     ordermodel.Currency
}
