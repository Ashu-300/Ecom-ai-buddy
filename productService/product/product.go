package product

import (
	"log"
	"supernova/productService/product/src/broker"
	"supernova/productService/product/src/db"
	"supernova/productService/product/src/routes"
	"supernova/productService/product/src/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
func SetupProductApp(router *gin.Engine) {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment")
	}

	db.InItDB()
	services.CloudinaryInit()
	broker.Connect()
	
	routes.ProductRoutes(router)
}

