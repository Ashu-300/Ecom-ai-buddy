package middlewares

import (
	"net/http"
	"strings"
	"supernova/authService/auth/src/db"
	"supernova/authService/auth/src/jwtutils"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		token := parts[1]

		claims, err := jwtutils.VerifyToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		expTime, err := claims.RegisteredClaims.GetExpirationTime()
		if err != nil || expTime.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}

		isBlackListed, err := db.IsTokenBlacklisted(token)

		if isBlackListed {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token has expired login again",
			})
			return
		}

		remainingTime := time.Until(expTime.Time)

		c.Set("remainingTime", remainingTime)
		c.Set("Email", claims.Email)
		c.Set("_id", claims.UserID)
		c.Set("token", token)

		c.Next()
	}
}
