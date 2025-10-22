package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"supernova/paymentService/payment/src/db"
	"supernova/paymentService/payment/src/dto"
	paymentmodel "supernova/paymentService/payment/src/paymentModel"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePayment(c *gin.Context) {
	orderServiceUrl := os.Getenv("ORDER_SERVICE_URL")
	orderID := c.Param("orderID")
	orderObjectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID format"})
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
	url := fmt.Sprintf("%s/api/order/get/%s", orderServiceUrl, orderID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to connect to order service"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	client := http.Client{Timeout: 10 * time.Second}
	orderResp, err := client.Do(req)
	if err != nil {
		c.JSON(orderResp.StatusCode, gin.H{
			"error": fmt.Sprintf("Cart service failed with status: %d", orderResp.StatusCode),
		})
		return
	}
	defer orderResp.Body.Close()

	if orderResp.StatusCode != http.StatusOK {
		c.JSON(orderResp.StatusCode, gin.H{
			"error": fmt.Sprintf("Order service failed with status: %d", orderResp.StatusCode),
		})
		return
	}
	var OrderDTO dto.OrderResponseDTO
	err = json.NewDecoder(orderResp.Body).Decode(&OrderDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode order details"})
		return
	}
	userOrder := OrderDTO.Order
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}
	userObjectID, _ := primitive.ObjectIDFromHex(userID.(string))
	paymentCollection := db.GetPaymentCollection()
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	filter := bson.M{
		"userID": userObjectID,
		"orderID":orderObjectID,
	}
	var existingPayment paymentmodel.Payment
	err = paymentCollection.FindOne(ctx , filter).Decode(&existingPayment)
	if err == nil {
		c.JSON(http.StatusBadRequest , gin.H{
			"error":"payment is already initiated",
			"paymentID":existingPayment.PaymentID,
		})
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var payment paymentmodel.Payment
	payment.PaymentID = primitive.NewObjectID()
	payment.OrderID = orderObjectID
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	payment.Price = paymentmodel.Price{
		TotalAmount: userOrder.TotalPrice.Amount,
		Currency:    paymentmodel.Currency(userOrder.TotalPrice.Currency),
	}
	payment.Status = paymentmodel.StatusPending
	payment.UserID = userObjectID

	result, err := paymentCollection.InsertOne(ctx, payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order in database", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":   "Payment initiated",
		"paymentID": result.InsertedID,
	})

}
