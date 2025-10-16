package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"supernova/orderService/order/src/db"
	"supernova/orderService/order/src/dto"
	"supernova/orderService/order/src/orderModel"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // Needed for OrderID
	"go.mongodb.org/mongo-driver/mongo"
)

// isValidStatusTransition checks if a status transition is valid


func CreateOrder(c *gin.Context) {
	// 1. Get Context Values and Type Assertions
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userToken, exists := c.Get("Token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Token not found in context"})
		return
	}
	tokenStr, ok := userToken.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error: Invalid token type"})
		return
	}

	client := http.Client{Timeout: 10 * time.Second}
	cartServiceURL := os.Getenv("CART_SERVICE_URL")
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")

	// ----------------------------------------------------
	// 2. Fetch User Cart (Product Service)
	// ----------------------------------------------------
	req, err := http.NewRequest("GET", cartServiceURL+"/api/cart/get", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	// Use a new variable for the response (cartResp)
	cartResp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to connect to cart service"})
		return
	}
	// Defer close immediately after a successful response
	defer cartResp.Body.Close()

	if cartResp.StatusCode != http.StatusOK {
		c.JSON(cartResp.StatusCode, gin.H{
			"error": fmt.Sprintf("Cart service failed with status: %d", cartResp.StatusCode),
		})
		return
	}

	var respCart dto.ResponseCart
	var userCart dto.Cart
	err = json.NewDecoder(cartResp.Body).Decode(&respCart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode cart details"})
		return
	}
	userCart = respCart.Cart

	if len(userCart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty. Cannot create order."})
		return
	}

	// ----------------------------------------------------
	// 3. Fetch User Details (Auth Service)
	// ----------------------------------------------------
	req, err = http.NewRequest("GET", authServiceURL+"/api/auth/user", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auth request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	// Use a new variable for the response (userResp)
	userResp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to connect to auth service"})
		return
	}
	// Defer close immediately after a successful response
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		c.JSON(userResp.StatusCode, gin.H{
			"error":         fmt.Sprintf("Auth service failed with status: %d", userResp.StatusCode),
			"response_body": userResp.Body,
		})
		return
	}

	var authResp dto.AuthResponse

	err = json.NewDecoder(userResp.Body).Decode(&authResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user details"})
		return
	}
	user := authResp.UserInfo

	// Check if user has at least one address
	if len(user.Addresses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User must have a shipping address configured."})
		return
	}

	var address dto.Address
	address = user.Addresses[0]

	// ----------------------------------------------------
	// 4. Calculate Total and Build Order Model
	// ----------------------------------------------------

	var totalAmount float64
	orderItems := make([]ordermodel.Item, 0, len(userCart.Items)) // preallocate slice

	for _, item := range userCart.Items {
		// Calculate total
		totalAmount += item.Price.Amount * float64(item.Quantity)

		// Convert dto.Item → ordermodel.Item and append
		orderItems = append(orderItems, ordermodel.Item{
			ProductID: item.ProductID,
			Price: ordermodel.Price{
				Amount:   item.Price.Amount,
				Currency: ordermodel.Currency(item.Price.Currency),
			},
			Quantity: item.Quantity,
		})
	}

	var order ordermodel.Order
	// Assign a new ObjectID here, as it's the MongoDB primary key
	order.OrderID = primitive.NewObjectID()
	order.UserID = userObjectID
	order.Items = orderItems // Spread operator to convert slice types
	order.TotalAmount = totalAmount
	order.Status = ordermodel.StatusPending // Directly assign constant
	order.Address = ordermodel.Address(address)
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// ----------------------------------------------------
	// 5. Save to Database
	// ----------------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderCollection := db.GetOrderCollection()
	result, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order in database", "details": err.Error()})
		return
	}

	cartClearReq, err := http.NewRequest("DELETE", cartServiceURL+"/api/cart/clear", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart clear request"})
		return
	}
	cartClearReq.Header.Set("Authorization", "Bearer "+tokenStr)

	cartClearResp, err := client.Do(cartClearReq)
	if err != nil || cartClearResp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to clear cart after order creation"})
		return
	}
	defer cartClearResp.Body.Close()

	// Final success response
	c.JSON(http.StatusCreated, gin.H{
		"message":       "Order created successfully",
		"order_details": order,
		"inserted_id":   result.InsertedID,
	})
}

func GetOrders(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	orderCollection := db.GetOrderCollection()

	filter := bson.M{"userId": userObjectID}
	cursor, err := orderCollection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{
				"message": "No orders found for the user",
				"orders":  []ordermodel.Order{},
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders from database", "details": err.Error()})
			return
		}

	}
	defer cursor.Close(ctx)

	var orders []ordermodel.Order
	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse orders from database",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Orders fetched successfully",
		"orders":  orders,
	})
}

func GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID format"})
		return
	}
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	orderCollection := db.GetOrderCollection()
	filter := bson.M{"_id": orderObjectID, "userId": userObjectID}
	var order ordermodel.Order
	err = orderCollection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order from database", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Order fetched successfully",
		"order":   order,
	})
}

func CancleOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID format"})
		return
	}
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	orderCollection := db.GetOrderCollection()
	filter := bson.M{"_id": orderObjectID, "userId": userObjectID}
	update := bson.M{
		"$set": bson.M{
			"status":    ordermodel.StatusCancelled,
			"updatedAt": time.Now(),
		},
	}
	var orderDoc ordermodel.Order
	err = orderCollection.FindOne(ctx, filter).Decode(&orderDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order from database", "details": err.Error()})
		return
	}

	if orderDoc.Status == ordermodel.StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is already cancelled"})
		return
	}
	if orderDoc.Status != ordermodel.StatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	result, err := orderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order in database", "details": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
	})
}

func UpdateOrderAddress(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID format"})
		return
	}
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	var addressDTO dto.Address
	if err := c.ShouldBindJSON(&addressDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address data", "details": err.Error()})
		return
	}
	update := bson.M{
		"$set": bson.M{
			"address":   ordermodel.Address(addressDTO),
			"updatedAt": time.Now(),
		},
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()
	orderCollection := db.GetOrderCollection()

	filter := bson.M{"userId": userObjectID, "_id": orderObjectID}

	result, err := orderCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update address",
			"details": err.Error(),
		})
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Address updated successfully",
	})

}

func isValidStatusTransition(currentStatus, newStatus ordermodel.OrderStatus) bool {
	// Define valid transitions
	validTransitions := map[ordermodel.OrderStatus][]ordermodel.OrderStatus{
		ordermodel.StatusPending: {
			ordermodel.StatusConfirmed,
			ordermodel.StatusCancelled,
		},
		ordermodel.StatusConfirmed: {
			ordermodel.StatusShipped,
			ordermodel.StatusCancelled,
		},
		ordermodel.StatusShipped: {
			ordermodel.StatusDelivered,
		},
		ordermodel.StatusDelivered: {}, // No valid transitions from delivered
		ordermodel.StatusCancelled: {}, // No valid transitions from cancelled
	}

	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, validStatus := range allowedTransitions {
		if validStatus == newStatus {
			return true
		}
	}
	return false
}

// UpdateOrderStatus updates the status of an existing order
func UpdateOrderStatus(c *gin.Context) {
	// Get the order ID from the URL parameter
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	// Convert string ID to ObjectID
	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	// Get the new status from the request body
	var updateRequest struct {
		Status ordermodel.OrderStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate the status value
	validStatuses := map[ordermodel.OrderStatus]bool{
		ordermodel.StatusPending:   true,
		ordermodel.StatusConfirmed: true,
		ordermodel.StatusShipped:   true,
		ordermodel.StatusDelivered: true,
		ordermodel.StatusCancelled: true,
	}

	if !validStatuses[updateRequest.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order status",
			"validStatuses": []ordermodel.OrderStatus{
				ordermodel.StatusPending,
				ordermodel.StatusConfirmed,
				ordermodel.StatusShipped,
				ordermodel.StatusDelivered,
				ordermodel.StatusCancelled,
			},
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	ctx := context.Background()
	collection := db.GetOrderCollection()

	// First, get the current order to check its status
	var currentOrder ordermodel.Order
	err = collection.FindOne(ctx, bson.M{"_id": orderObjectID}).Decode(&currentOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order details"})
		return
	}

	isAdminVal, exists := c.Get("IsAdmin")
	isAdmin, ok := isAdminVal.(bool)

	if !exists || !ok || !isAdmin {
		// not admin — so only allow if user owns the order
		if currentOrder.UserID != userObjectID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this order"})
			return
		}
	}

	

	// Validate status transition
	if !isValidStatusTransition(currentOrder.Status, updateRequest.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid status transition from %s to %s",
				currentOrder.Status, updateRequest.Status),
		})
		return
	}

	// Create update document
	update := bson.M{
		"$set": bson.M{
			"status":    updateRequest.Status,
			"updatedAt": time.Now(),
		},
	}

	// Update the order
	filter := bson.M{"_id": orderObjectID}
	if !isAdmin {
		filter["userId"] = userObjectID 
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Order status updated successfully",
		"orderId":   orderID,
		"oldStatus": currentOrder.Status,
		"newStatus": updateRequest.Status,
	})
}
