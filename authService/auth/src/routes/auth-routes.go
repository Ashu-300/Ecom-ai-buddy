package routes

import (
	"supernova/authService/auth/src/controller"
	"supernova/authService/auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	publicRoutes := router.Group("/api/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login" , controller.Login)

	securedRoutes := router.Group("/api/auth/")
	
	securedRoutes.Use(middlewares.AuthMiddleware()) 
	securedRoutes.GET("/user" , controller.GetCurrentUser)
	securedRoutes.POST("/logout" , controller.Logout)


}