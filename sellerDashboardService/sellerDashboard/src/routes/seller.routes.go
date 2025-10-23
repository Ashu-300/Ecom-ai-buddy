package routes

import (
	"supernova/sellerDashboardService/sellerDashboard/src/controller"
	"supernova/sellerDashboardService/sellerDashboard/src/middleware"

	"github.com/gin-gonic/gin"
)

func SellerRoutes(router *gin.Engine){
	r := router.Group("/api/sellerdashboard")

	securedRoutes := r.Use(middleware.CreateAuthMiddleware())

	securedRoutes.GET("/get/metrics" , controller.GetMetrics)
	securedRoutes.GET("/get/order" , controller.GetOrders)
	securedRoutes.GET("/get/product" , controller.GetProducts)
}