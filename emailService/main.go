package main

import (
	"log"
	"supernova/emailService/email"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
	
	router := gin.Default()

	email.SetupEmailApp(router)

	router.Run(":8087")
}