package dto

type Item struct {
	ProductID string 				`bson:"productId" json:"productId"`
	Price     Price              	`bson:"price" json:"price"`
	Quantity  int                	`bson:"quantity" json:"quantity" binding:"required,min=1"`
}
type Price struct {
	Amount  	float64 `bson:"amount" json:"amount"`
	Currency	Currency  `bson:"currency" json:"currency"`
}
type Currency string

const (
	USD Currency = "USD"
	INR Currency = "INR"
)