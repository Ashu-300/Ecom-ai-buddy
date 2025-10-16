package cartroutes

import (
	"supernova/cartService/cart/src/cartController"
	"supernova/cartService/cart/src/cartMiddleware"

	"github.com/gin-gonic/gin"
)

func SetupCartRoutes(router *gin.Engine) {

	r := router.Group("/api/cart")

	securedRoutes := r.Use(cartmiddleware.CreateAuthMiddleware())

	securedRoutes.POST("/item" , cartcontroller.AddItemToCart )
	securedRoutes.PATCH("/updateitem" , cartcontroller.UpdateItemQuantity)
	securedRoutes.PATCH("removeitem" , cartcontroller.RemoveItemFromCart)
	securedRoutes.GET("/get" , cartcontroller.GetCart)
	securedRoutes.DELETE("/clear" , cartcontroller.ClearCart)
}