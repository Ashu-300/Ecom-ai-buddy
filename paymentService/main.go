package main

import (
	"supernova/paymentService/payment"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	payment.SetupPaymentApp(router)

	router.Run(":8085")
}