package main

import (
	"supernova/sellerDashboardService/sellerDashboard"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	sellerdashboard.SetupSellerDashboardApp(router)

	router.Run(":8088")
}