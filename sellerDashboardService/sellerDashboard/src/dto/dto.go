package dto

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
 	"github.com/golang-jwt/jwt/v5"

)

type Claims struct {
	Email  string `json:"username"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

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


