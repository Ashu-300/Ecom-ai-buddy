package dto

import (
	cartmodel "supernova/cartService/cart/src/cartModel"
)

type Item struct {
	ProductID string 			`bson:"productId" json:"productId"`
	Price     cartmodel.Price  `bson:"price" json:"price"`
	Quantity  int    			`bson:"quantity" json:"quantity" binding:"required,min=1"`
}
