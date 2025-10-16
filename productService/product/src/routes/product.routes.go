package routes

import (
	"supernova/productService/product/src/controllers"
	"supernova/productService/product/src/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	r := router.Group("/api/product")

	r.GET("/get",controllers.GetProducts)
	r.GET("/get/:id",controllers.GetProductByID)

	securedRoute := r.Use(middleware.CreateAuthMiddleware())
	
	securedRoute.POST("/create",controllers.CreateProduct)
	securedRoute.PATCH("/:id" ,controllers.UpdateProduct)

	

}