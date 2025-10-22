package dto

import (
	"time"
)

// OrderResponseDTO represents the top-level structure of the JSON response.
type OrderResponseDTO struct {
	Message string `json:"message"`
	Order   OrderDTO `json:"order"`
}

// OrderDTO represents the main order details.
type OrderDTO struct {
	OrderID   string       `json:"orderId"`
	UserID    string       `json:"userId"`
	Items     []ItemDTO    `json:"items"`
	TotalPrice PriceDTO    `json:"totalPrice"` // Use float64 or int, depending on how amount is stored
	Status    string       `json:"status"`
	Address   AddressDTO   `json:"address"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// ItemDTO represents a single product item within the order.
type ItemDTO struct {
	ProductID string    `json:"productId"`
	Price     PriceDTO  `json:"price"`
	Quantity  int       `json:"quantity"`
}

// PriceDTO represents the price structure for an item.
type PriceDTO struct {
	Amount   float64 `json:"amount"` // Use float64 for currency precision
	Currency string  `json:"currency"`
}

// AddressDTO represents the shipping or billing address.
type AddressDTO struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}