package dto

import (
    "mime/multipart"
)

// PriceDTO for receiving price info
type PriceDTO struct {
    Amount   float64 `form:"amount" binding:"required"`
    Currency string  `form:"currency" binding:"required,oneof=USD INR"`
}

// ImageDTO for receiving image files
type ImageDTO struct {
    File      *multipart.FileHeader `form:"file" binding:"required"` // This will only accept one file at a time
}

// ProductDTO for receiving product data along with images
type ProductDTO struct {
    Title       string       `form:"title" binding:"required"`
    Description string       `form:"description"`
    Price       PriceDTO     `form:"price" binding:"required"`
    Images      []*multipart.FileHeader `form:"images" binding:"required"`
    Stock       int          `form:"stock" binding:"required,gte=0"`
    // SellerID    string       `form:"seller_id" binding:"required"`
}
