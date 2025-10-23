package sellerdashboard

import (
	"log"
	"supernova/sellerDashboardService/sellerDashboard/src/broker"
	"supernova/sellerDashboardService/sellerDashboard/src/db"
	"supernova/sellerDashboardService/sellerDashboard/src/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupSellerDashboardApp(router *gin.Engine){
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment")
	}

	db.InitDB()
	broker.Connect()
	broker.ConsumeQueues()
	routes.SellerRoutes(router)

}