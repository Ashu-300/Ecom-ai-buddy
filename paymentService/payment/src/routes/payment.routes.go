package routes

import (
	"supernova/paymentService/payment/src/controller"
	"supernova/paymentService/payment/src/middleware"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine){
	r := router.Group("/api/payment")

	securedRoutes := r.Use(middleware.CreateAuthMiddleware())

	securedRoutes.POST("/create/:orderID" , controller.CreatePayment)
	securedRoutes.POST("/verify/:paymentID" , controller.VerifyPayment)
}