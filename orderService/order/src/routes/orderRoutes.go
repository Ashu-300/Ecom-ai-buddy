package routes

import (
	"supernova/orderService/order/src/controller"
	"supernova/orderService/order/src/middleware"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine) {
	r := router.Group("/api/order")

	securedRoutes := r.Use(middleware.CreateAuthMiddleware())

	securedRoutes.POST("/create" , controller.CreateOrder)
	securedRoutes.GET("/get" , controller.GetOrders)
	securedRoutes.GET("/get/:id" , controller.GetOrderByID)
	securedRoutes.PATCH("/cancle/:id" , controller.CancleOrderByID)
	securedRoutes.PATCH("/update/address/:id" , controller.UpdateOrderAddress)
	securedRoutes.PATCH("/update/status/:id" , controller.UpdateOrderStatus)

}