package main

import (
	"strings"
	"time"

	"supernova/authService/auth"
	"supernova/cartService/cart"
	"supernova/emailService/email"
	"supernova/orderService/order"
	"supernova/paymentService/payment"
	"supernova/productService/product"
	sellerdashboard "supernova/sellerDashboardService/sellerDashboard"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/prometheus/client_golang/prometheus"
)

// ---------------- METRICS ----------------

var httpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Service level latency",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"service"},
)

func initMetrics() {
	prometheus.MustRegister(httpRequestDuration)
}

// ---------------- MIDDLEWARE ----------------

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		path := c.FullPath() // IMPORTANT

		service := extractService(path)

		httpRequestDuration.WithLabelValues(service).Observe(duration)
	}
}

// ---------------- SERVICE EXTRACTION ----------------

func extractService(path string) string {
	switch {
	case strings.HasPrefix(path, "/auth"):
		return "auth"
	case strings.HasPrefix(path, "/cart"):
		return "cart"
	case strings.HasPrefix(path, "/product"):
		return "product"
	case strings.HasPrefix(path, "/order"):
		return "order"
	case strings.HasPrefix(path, "/payment"):
		return "payment"
	case strings.HasPrefix(path, "/email"):
		return "email"
	case strings.HasPrefix(path, "/seller"):
		return "seller"
	default:
		return "unknown"
	}
}

// ---------------- MAIN ----------------

func main() {
	router := gin.Default()

	// 1️⃣ default gin metrics (optional but useful)
	prom := ginprometheus.NewPrometheus("gin")
	prom.Use(router)

	// 2️⃣ custom service-level metrics
	initMetrics()
	router.Use(prometheusMiddleware())

	// 3️⃣ register services
	auth.SetupAuthApp(router)
	cart.SetupCartApp(router)
	product.SetupProductApp(router)
	order.SetupOrderApp(router)
	payment.SetupPaymentApp(router)
	email.SetupEmailApp(router)
	sellerdashboard.SetupSellerDashboardApp(router)

	// 4️⃣ run server
	router.Run(":8080")
}