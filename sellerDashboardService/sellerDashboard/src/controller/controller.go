package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"supernova/sellerDashboardService/sellerDashboard/src/db"
	"supernova/sellerDashboardService/sellerDashboard/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(user models.User) {
	userCollection := db.GetSellerUserCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := userCollection.InsertOne(ctx , user)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}

func CreateProduct(product models.Product){
	productCollection := db.GetSellerProductCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := productCollection.InsertOne(ctx , product)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}


func CreateOrder(order models.Order){
	orderCollection := db.GetSellerOrderCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := orderCollection.InsertOne(ctx , order)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}

func CreatePayment(payment models.Order){
	paymentCollection := db.GetSellerPaymentCollection()
	ctx , cancle := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancle()
	_ , err := paymentCollection.InsertOne(ctx , payment)
	if err != nil {
		log.Print("error: %v",err.Error())
	}
}


// func GetMetrics(c *gin.Context) {
// 	// sellerID, exists := c.Get("UserID")
// 	// if !exists {
// 	// 	log.Println("❌ Seller ID not found in context")
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized: seller id missing"})
// 	// 	return
// 	// }
// 	// sellerIDStr := fmt.Sprintf("%v", sellerID)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// --- Current month date range ---
// 	now := time.Now()
// 	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
// 	// Optional: last day is not strictly needed with $gte and $lt next month
// 	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

// 	filter := bson.M{
// 		"status": "delivered",
// 		"createdAt": bson.M{
// 			"$gte": firstOfMonth,
// 			"$lt":  firstOfNextMonth,
// 		},
// 	}

// 	cursor, err := db.GetSellerOrderCollection().Find(ctx, filter)
// 	if err != nil {
// 		log.Println("❌ Error fetching orders:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
// 		return
// 	}
// 	log.Print(cursor)
// 	defer cursor.Close(ctx)

// 	var orders []models.Order
// 	if err := cursor.All(ctx, &orders); err != nil {
// 		log.Println("❌ Error decoding orders:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode orders"})
// 		return
// 	}

// 	// --- Calculate total sales count and total revenue ---
// 	totalRevenue := 0.0
// 	totalSalesCount := 0
// 	for _, order := range orders {
// 		for _, item := range order.Items {
// 			totalRevenue += item.Price.Amount * float64(item.Quantity)
// 			totalSalesCount += item.Quantity
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"totalRevenue":    totalRevenue,
// 		"totalSalesCount": totalSalesCount,
// 		"ordersCount":     len(orders),
// 	})
// }



type TopProduct struct {
	ProductID primitive.ObjectID `json:"productId"`
	SoldUnits int                `json:"soldUnits"`
	Revenue   float64            `json:"revenue"`
}

func GetMetrics(c *gin.Context) {
	

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- Current month filter ---
	now := time.Now().UTC()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	filter := bson.M{
		"status": "delivered",
		"createdAt": bson.M{
			"$gte": firstOfMonth,
			"$lt":  firstOfNextMonth,
		},
	}

	cursor, err := db.GetSellerOrderCollection().Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode orders"})
		return
	}

	// --- Calculate totals and top product ---
	totalRevenue := 0.0
	totalSalesCount := 0

	// Map[productID] -> TopProduct
	productStats := make(map[primitive.ObjectID]*TopProduct)

	for _, order := range orders {
		for _, item := range order.Items {
			// Only consider items belonging to this seller (optional: if Product info available)
			totalRevenue += item.Price.Amount * float64(item.Quantity)
			totalSalesCount += item.Quantity

			if tp, ok := productStats[item.ProductID]; ok {
				tp.SoldUnits += item.Quantity
				tp.Revenue += item.Price.Amount * float64(item.Quantity)
			} else {
				productStats[item.ProductID] = &TopProduct{
					ProductID: item.ProductID,
					SoldUnits: item.Quantity,
					Revenue:   item.Price.Amount * float64(item.Quantity),
				}
			}
		}
	}

	// Find top product
	var topProduct *TopProduct
	for _, tp := range productStats {
		if topProduct == nil || tp.SoldUnits > topProduct.SoldUnits {
			topProduct = tp
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"totalRevenue":    totalRevenue,
		"totalSalesCount": totalSalesCount,
		"ordersCount":     len(orders),
		"topProduct":      topProduct,
	})
}


func GetOrders(c *gin.Context) {
	sellerID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "seller ID missing"})
		return
	}
	sellerIDStr := fmt.Sprintf("%v", sellerID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- 1️⃣ Get all products for this seller ---
	productCursor, err := db.GetSellerProductCollection().Find(ctx, bson.M{"seller_id": sellerIDStr})
	if err != nil {
		log.Println("❌ Error fetching products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}
	defer productCursor.Close(ctx)

	var products []models.Product
	if err := productCursor.All(ctx, &products); err != nil {
		log.Println("❌ Error decoding products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode products"})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, []models.Order{}) // no products, return empty list
		return
	}

	productIDs := make([]primitive.ObjectID, 0, len(products))
	for _, p := range products {
		id, err := primitive.ObjectIDFromHex(p.SellerID) // adjust if your Product._id is ObjectID
		if err == nil {
			productIDs = append(productIDs, id)
		}
	}

	// --- 2️⃣ Get all orders containing seller's products ---
	orderFilter := bson.M{
		"items.productId": bson.M{"$in": productIDs},
	}

	orderCursor, err := db.GetSellerOrderCollection().Find(ctx, orderFilter)
	if err != nil {
		log.Println("❌ Error fetching orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}
	defer orderCursor.Close(ctx)

	var orders []models.Order
	if err := orderCursor.All(ctx, &orders); err != nil {
		log.Println("❌ Error decoding orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode orders"})
		return
	}

	// --- 3️⃣ Filter order items to only include those from this seller ---
	filteredOrders := make([]models.Order, 0)
	for _, order := range orders {
		filteredItems := make([]models.Item, 0)
		for _, item := range order.Items {
			for _, pid := range productIDs {
				if item.ProductID == pid {
					filteredItems = append(filteredItems, item)
					break
				}
			}
		}
		if len(filteredItems) > 0 {
			order.Items = filteredItems
			filteredOrders = append(filteredOrders, order)
		}
	}

	c.JSON(http.StatusOK, filteredOrders)
}


func GetProducts(c *gin.Context) {
	sellerID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "seller ID missing"})
		return
	}
	sellerIDStr := fmt.Sprintf("%v", sellerID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all products for this seller, sorted by createdAt descending
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := db.GetSellerProductCollection().Find(ctx, bson.M{"seller_id": sellerIDStr}, findOptions)
	if err != nil {
		log.Println("❌ Error fetching products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		log.Println("❌ Error decoding products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
