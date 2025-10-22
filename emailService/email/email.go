package email

import (
	"supernova/emailService/email/broaker"

	"github.com/gin-gonic/gin"
)

func SetupEmailApp(router *gin.Engine){

	broaker.Connect()
}