package main

import (
	"supernova/authService/auth"
	"supernova/cartService/cart"
	"supernova/orderService/order"
	"supernova/paymentService/payment"
	"supernova/productService/product"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	auth.SetupAuthApp(router)     	// port 8081
	cart.SetupCartApp(router) 		// port 8082
	product.SetupProductApp(router) // port 8083
	order.SetupOrderApp(router) 	// port 8084
	payment.SetupPaymentApp(router)	// port 8085

	router.Run(":8080")
}