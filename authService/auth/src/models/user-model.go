package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type Address struct {
	Street    	string `json:"street" binding:"required"`
	City      	string `json:"city" binding:"required"`
	State     	string `json:"state" binding:"required"`
	PostalCode	string `json:"postal_code" `
	Country   	string `json:"country" binding:"required"`
}

type User struct {
	ID			primitive.ObjectID	 `bson:"_id,omitempty" json:"id"`
    UserName  	string `json:"username" binding:"required"`
    Email     	string `json:"email" binding:"required,email"`
    Password  	string `json:"password" binding:"required,min=6"`
    FirstName 	string `json:"first_name" binding:"required"`
    LastName  	string `json:"last_name" binding:"required"`
	Role 	  	string `json:"role" binding:"oneof=seller user"`
	Addresses 	[]Address `json:"addresses" binding:"dive"`
}


