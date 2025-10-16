package dto

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email  string `json:"username"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}