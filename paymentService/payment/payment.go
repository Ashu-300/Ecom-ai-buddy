package payment

import (
	"supernova/paymentService/payment/src/db"
	"supernova/paymentService/payment/src/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupPaymentApp(router *gin.Engine) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db.InitDB()

	routes.PaymentRoutes(router)

}