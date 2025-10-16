package order

import (
	"log"
	"supernova/orderService/order/src/db"
	"supernova/orderService/order/src/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupOrderApp(router *gin.Engine) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db.InitDB()
	log.Print("order service")

	routes.SetupOrderRoutes(router)

}