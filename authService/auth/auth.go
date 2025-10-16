package auth

import (
	"log"
	"supernova/authService/auth/src/db"
	"supernova/authService/auth/src/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupAuthApp(router *gin.Engine) {

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment")
	}

	// Init DBs
	db.InitDB()
	db.InitRedisDB()
	db.CreateUserIndexes(db.UserCollection)

	// Setup router

	routes.AuthRoutes(router)

}
