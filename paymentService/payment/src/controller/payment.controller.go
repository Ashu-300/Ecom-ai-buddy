package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"supernova/paymentService/payment/src/broker"
	"supernova/paymentService/payment/src/db"
	"supernova/paymentService/payment/src/dto"
	paymentmodel "supernova/paymentService/payment/src/paymentModel"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
type JsonPayment struct{
	ReceiverMail string 			`json:"receiverMail"`
	PaymentID  	primitive.ObjectID `json:"paymentID"`
	OrderID		primitive.ObjectID `json:"orderID"`
	Amount 		float64				`json:"amount"`
	Currency    string				`json:"currency"`
}

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
	userEmail , exists := c.Get("Email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: email not found in context"})
		return
	}
	emailStr , ok := userEmail.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error: Invalid email type"})
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
	Jsonpayment := JsonPayment{
		ReceiverMail: emailStr,
		OrderID: payment.OrderID,
		PaymentID: payment.PaymentID,
		Amount : payment.Price.TotalAmount,
		Currency: userOrder.TotalPrice.Currency,
	}
	body , err := json.Marshal(Jsonpayment)
	if err != nil {
        log.Printf("Failed to marshal user: %v", err)
    }

	err = broker.PublishJSON("PaymentService" , body )
	if err != nil {
        log.Print("Error in sending message to broaker" , err.Error())
    }

	paymentJson , err := json.Marshal(payment)
	if err != nil {
        log.Printf("Failed to marshal user: %v", err)
    }
	err = broker.PublishJSON("PaymentDashboard" , paymentJson )
	if err != nil {
        log.Print("Error in sending message to broaker" , err.Error())
    }

	c.JSON(200, gin.H{
		"message":   "Payment initiated",
		"paymentID": result.InsertedID,
	})

}


func VerifyPayment(c *gin.Context) {
    paymentCollection := db.GetPaymentCollection()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // --- Step 1: Extract payment ID from URL ---
    paymentID := c.Param("paymentID")
    paymentObjectID, err := primitive.ObjectIDFromHex(paymentID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
        return
    }

    // --- Step 2: Parse body for verification data (e.g., transaction ID, status) ---
    var verifyReq struct {
        PaymentGatewayTxnID string `json:"paymentGatewayTxnID"`
        Status              string `json:"status"` // "success" / "failed"
    }
    if err := c.ShouldBindJSON(&verifyReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // --- Step 3: Fetch payment from DB ---
    var existingPayment paymentmodel.Payment
    err = paymentCollection.FindOne(ctx, bson.M{"_id": paymentObjectID}).Decode(&existingPayment)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // --- Step 4: Prevent double verification ---
    if existingPayment.Status == paymentmodel.StatusCompleted {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Payment already verified as successful"})
        return
    }

    // --- Step 5: Update status based on verification result ---
    var newStatus paymentmodel.Status
    if verifyReq.Status == "success" {
        newStatus = paymentmodel.StatusCompleted
    } else {
        newStatus = paymentmodel.StatusFailed
    }

    update := bson.M{
        "$set": bson.M{
            "status":     newStatus,
            "txnID":      verifyReq.PaymentGatewayTxnID,
            "updatedAt":  time.Now(),
        },
    }

    _, err = paymentCollection.UpdateOne(ctx, bson.M{"_id": paymentObjectID}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
        return
    }

    
    
    // --- Step 7: Respond ---
    c.JSON(http.StatusOK, gin.H{
        "message":    "Payment verification updated",
        "status":     newStatus,
        "paymentID":  paymentID,
        "orderID":    existingPayment.OrderID,
    })
}
