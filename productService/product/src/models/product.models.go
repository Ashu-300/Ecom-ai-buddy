package models


// import "go.mongodb.org/mongo-driver/bson/primitive"

// Price sub-struct
type Price struct {
    Amount   float64  `bson:"amount" json:"amount" binding:"required"`
    Currency string  `bson:"currency" json:"currency" binding:"required,oneof=USD INR"`
}

// Image sub-struct
type Image struct {
    URL       string `bson:"url" json:"url"`
    Thumbnail string `bson:"thumbnail" json:"thumbnail"`
	ID 	  	  string `bson:"id" json:"id"`
}

// Product model
type Product struct {
    Title       string             `bson:"title" json:"title" binding:"required"`
    Description string             `bson:"description" json:"description"`
    Price       Price              `bson:"price" json:"price"  binding:"required"`
    Images      []Image            `bson:"images" json:"images" binding:"required"`
    Stock       int                `bson:"stock" json:"stock" binding:"required,gte=0"`
    SellerID    string             `bson:"seller_id" json:"seller_id" binding:"required"`
}

