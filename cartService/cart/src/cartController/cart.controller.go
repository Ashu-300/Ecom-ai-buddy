package cartcontroller

import (
	"context"
	"net/http"
	"supernova/cartService/cart/src/cartModel"
	"supernova/cartService/cart/src/db"
	"supernova/cartService/cart/src/dto"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddItemToCart(c *gin.Context) {
	var item dto.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, gin.H{"error": "UserID not found in context"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	productObjectID, err := primitive.ObjectIDFromHex(item.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}
	item.ProductID = productObjectID.Hex()

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	res := db.GetCartCollection().FindOne(ctx, bson.M{"userId": userObjectID})
	cart, err := res.Raw()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newCart := cartmodel.Cart{}
			newCart.UserID = userObjectID
			newCart.Items = []cartmodel.Item{
				{
					productObjectID,
					item.Price,
					item.Quantity,
				},
			}
			newCart.CreatedAt = time.Now()
			newCart.UpdatedAt = time.Now()
			_, err := db.GetCartCollection().InsertOne(ctx, newCart)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to create cart"})
				return
			}
			c.JSON(200, gin.H{
				"message": "Item added to cart",
				"cart":    newCart,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	var existingCart cartmodel.Cart
	err = bson.Unmarshal(cart, &existingCart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
		return
	}

	itemExists := false
	for i, cartItem := range existingCart.Items {
		if cartItem.ProductID == productObjectID {
			existingCart.Items[i].Quantity += item.Quantity
			itemExists = true
			break
		}
	}
	if !itemExists {
		existingCart.Items = append(existingCart.Items, cartmodel.Item{
			ProductID: productObjectID,
			Quantity:  item.Quantity,
		})
	}
	existingCart.UpdatedAt = time.Now()

	_, err = db.GetCartCollection().UpdateOne(ctx, bson.M{"userId": userObjectID}, bson.M{"$set": existingCart})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update cart"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Item added to cart",
		"cart":    existingCart,
	})
	return

}

func UpdateItemQuantity(c *gin.Context) {
	var item cartmodel.Item

	err := c.ShouldBind(&item)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userId, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, gin.H{"error": "UserID not found in context"})
		return
	}
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	res := db.GetCartCollection().FindOne(ctx, bson.M{"userId": userId})
	cart, err := res.Raw()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var existingCart cartmodel.Cart
	err = bson.Unmarshal(cart, &existingCart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
		return
	}
	itemExists := false
	for i, cartItem := range existingCart.Items {
		if cartItem.ProductID == item.ProductID {
			existingCart.Items[i].Quantity = item.Quantity
			itemExists = true
			break
		}
	}
	if !itemExists {
		c.JSON(404, gin.H{"error": "Item not found in cart"})
		return
	}
	existingCart.UpdatedAt = time.Now()

	_, err = db.GetCartCollection().UpdateOne(ctx, bson.M{"userId": userId}, bson.M{"$set": existingCart})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update cart"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Item quantity updated",
		"cart":    existingCart,
	})
	return
}

func RemoveItemFromCart(c *gin.Context) {
	var item cartmodel.Item
	err := c.ShouldBind(&item)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userId, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, gin.H{"error": "UserID not found in context"})
		return
	}
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	res := db.GetCartCollection().FindOne(ctx, bson.M{"userId": userId})
	cart, err := res.Raw()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var existingCart cartmodel.Cart
	err = bson.Unmarshal(cart, &existingCart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
		return
	}
	itemIndex := -1
	for i, cartItem := range existingCart.Items {
		if cartItem.ProductID == item.ProductID {
			itemIndex = i
			break
		}
	}
	if itemIndex == -1 {
		c.JSON(404, gin.H{"error": "Item not found in cart"})
		return
	}
	existingCart.Items = append(existingCart.Items[:itemIndex], existingCart.Items[itemIndex+1:]...)
	existingCart.UpdatedAt = time.Now()
	_, err = db.GetCartCollection().UpdateOne(ctx, bson.M{"userId": userId}, bson.M{"$set": existingCart})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update cart"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Item removed from cart",
		"cart":    existingCart,
	})
	return
}

func GetCart(c *gin.Context) {
	userId, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, gin.H{"error": "UserID not found in context"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	res := db.GetCartCollection().FindOne(ctx, bson.M{"userId": userObjectID})
	cart, err := res.Raw()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var existingCart cartmodel.Cart
	err = bson.Unmarshal(cart, &existingCart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
		return
	}
	c.JSON(200, gin.H{
		"cart": existingCart,
	})
	return
}

func ClearCart(c *gin.Context) {
	userId, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, gin.H{"error": "UserID not found in context"})
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	_, err = db.GetCartCollection().DeleteOne(ctx, bson.M{"userId": userObjectID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to clear cart"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Cart cleared",
	})
	return
}
