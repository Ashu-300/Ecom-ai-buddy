package main

import (
	"supernova/cartService/cart"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	cart.SetupCartApp(router)

	router.Run(":8082")
}