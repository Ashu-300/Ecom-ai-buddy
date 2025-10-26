package main

import (
	"supernova/authService/auth"
	"supernova/cartService/cart"
	"supernova/emailService/email"
	"supernova/orderService/order"
	"supernova/paymentService/payment"
	"supernova/productService/product"
	"supernova/sellerDashboardService/sellerDashboard"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {
	router := gin.Default()

	prom := ginprometheus.NewPrometheus("gin")
	prom.Use(router)


	auth.SetupAuthApp(router)     					// port 8081
	cart.SetupCartApp(router) 						// port 8082
	product.SetupProductApp(router) 				// port 8083
	order.SetupOrderApp(router) 					// port 8084
	payment.SetupPaymentApp(router)					// port 8085
	// aibuddy.SetUpAibuddyAPP(router)				// port 8086
	email.SetupEmailApp(router)						// port 8087
	sellerdashboard.SetupSellerDashboardApp(router)	// port 8088

	router.Run(":8080")
}