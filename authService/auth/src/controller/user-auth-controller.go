package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"supernova/authService/auth/src/broker"
	"supernova/authService/auth/src/db"
	"supernova/authService/auth/src/dto"
	"supernova/authService/auth/src/jwtutils"
	"supernova/authService/auth/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)
type JsonUser struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}


func Register(c *gin.Context) {
    var newUser models.User

    // Validate input
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Set a new ObjectID for the user
    newUser.ID = primitive.NewObjectID()

    // Hash the password
    hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    newUser.Password = string(hashPassword)

    _, err = db.UserCollection.InsertOne(c, newUser)
    if err != nil {
        if mongo.IsDuplicateKeyError(err) {
            var existingUser models.User
            ctx, cancel := context.WithTimeout(c, 5*time.Second) 
            defer cancel()

            _ = db.UserCollection.FindOne(ctx, bson.M{"email": newUser.Email}).Decode(&existingUser)
            if existingUser.Email == newUser.Email {
                c.JSON(http.StatusConflict, gin.H{"error": "Email already exists."})
                return
            }
            

            c.JSON(http.StatusConflict, gin.H{"error": "A user with this email or username already exists."})
            return
        }

        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user."})
        return
    }

    jasonUser := JsonUser{
        Name:  newUser.FirstName,
        Email: newUser.Email,
    }
    body, err := json.Marshal(jasonUser)
    if err != nil {
        log.Printf("Failed to marshal user: %v", err)
    }

   

    // Attempt to publish a message to the queue.
    err = broker.PublishJSON("AuthService" , body)
    if err != nil {
        log.Print("Error in sending message to broaker" , err.Error())
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})
}

func Login(c *gin.Context) {
    var credentials dto.LoginCredential
    var user models.User

    // Bind JSON body
    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Filter to find user by email
    filter := bson.M{"email": credentials.Email}

    // Find user in MongoDB
    err := db.UserCollection.FindOne(c, filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // User not found
            c.JSON(http.StatusNotFound, gin.H{"message": "user not registered"})
            return
        }
        // Some other DB error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Check password
    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
        return
    }

    // Generate JWT token
    tokenString, err := jwtutils.GeneratejwtToken(user.ID.Hex(), user.Email , user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "could not generate token",
            "error":   err.Error(),
        })
        return
    }
    

       
    

    // Login successful
    c.JSON(http.StatusOK, gin.H{
        "message": "login success",
        "userId":  user.ID.Hex(),
        "email":   user.Email,
        "role":    user.Role,
        "token":   tokenString,
    })
}

func GetCurrentUser(c *gin.Context){
	var user dto.UserResponse ;

	 UserEmail , exists := c.Get("Email") ;
	 if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Email not found in context"}) ;
    	return
	 }
	 filter := bson.M{"email": UserEmail} ;

	  err := db.UserCollection.FindOne(c, filter).Decode(&user) ;

	  if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound , gin.H{
			"message" : "User Not Found",
		})
		return ;
	  }
	  
	  c.AbortWithStatusJSON(http.StatusOK , gin.H{
		"message" : "User Found successfully",
		"userInfo" : user ,
	  })
	  return ;
}

func Logout(c *gin.Context){
	tokenInterface , exists := c.Get("token") ;
	token := tokenInterface.(string)
	
	if !exists {
		c.JSON(http.StatusNotFound , gin.H{
			"message" : "token not found" ,
		})
		return ;
	}
	
	remainingTimeInterface , _ := c.Get("remainingTime")
	remainingTime := remainingTimeInterface.(time.Duration)
	
	if remainingTime <= 0 {
		c.JSON(http.StatusOK , gin.H{
			"message" : "token already expired",
		})
		return ;
	}

	err := db.BlacklistToken(token, remainingTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to blacklist token",
			"error":   err.Error(),
		})
		return ;
	}

	if err != nil {
		c.JSON(http.StatusConflict , gin.H{
			"message" : "can not blacklist the token",
		})
		return ;
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful! Token has been blacklisted.",
	})

	return ;

}