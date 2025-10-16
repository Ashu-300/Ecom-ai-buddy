package middleware

import (
	"net/http"
	"strings"
	"supernova/paymentService/payment/src/jwtutils"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
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

		// Verify token
		claims, err := jwtutils.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Check expiration
		expTime, err := claims.RegisteredClaims.GetExpirationTime()
		if err != nil || expTime.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}

		// Role check
		if claims.Role != "user"  {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		// Optional: check if blacklisted
		// if isBlacklisted(token) { ... }

		// Add values to context
		c.Set("remainingTime", time.Until(expTime.Time))
		c.Set("Email", claims.Email)
		c.Set("UserID", claims.UserID)
		c.Set("Token", token)
		c.Set("Role", claims.Role)

		c.Next()
	}
}
