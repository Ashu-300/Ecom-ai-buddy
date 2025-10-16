package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	Street     string `json:"street" `
	City       string `json:"city" `
	State      string `json:"state" `
	PostalCode string `json:"postal_code" `
	Country    string `json:"country" `
}

type UserResponse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName  string             `json:"username" `
	Email     string             `json:"email" `
	FirstName string   			 `json:"first_name" `
	LastName  string             `json:"last_name" `
	Role      string             `json:"role" `
	Addresses []Address          `json:"addresses`
}

type LoginCredential struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}