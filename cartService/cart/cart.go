package cart

import (
	"log"
	cartroutes "supernova/cartService/cart/src/cartRoutes"
	"supernova/cartService/cart/src/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupCartApp(router *gin.Engine) {

	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	db.InitDB()

	cartroutes.SetupCartRoutes(router)

}
