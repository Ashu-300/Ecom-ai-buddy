package main

import (
	"supernova/productService/product"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	product.SetupProductApp(router)

	router.Run(":8083")
}