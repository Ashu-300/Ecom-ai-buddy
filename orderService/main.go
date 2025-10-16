package main

import (
	"supernova/orderService/order"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	order.SetupOrderApp(router)

	router.Run(":8084")
}