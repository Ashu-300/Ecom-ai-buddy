package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupPaymentApp(router *gin.Engine) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

}