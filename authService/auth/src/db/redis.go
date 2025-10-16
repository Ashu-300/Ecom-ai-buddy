package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedisDB() {
	// Get Redis connection details from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	
	redisPort := os.Getenv("REDIS_PORT")
	
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Create the Redis client using the environment variables
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	ctx := context.Background()

	// Ping the Redis server to check the connection
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("✅ Auth service connected to Redis")
}

// BlacklistToken adds a token to the blacklist with a specific expiration.
func BlacklistToken(token string, expiration time.Duration) error {
    // The key to store the blacklisted token
    key := fmt.Sprintf("blacklist:%s", token)

    // Set the key in Redis with the token value and expiration
    err := rdb.Set(ctx, key, true, expiration).Err()
    if err != nil {
        return fmt.Errorf("failed to blacklist token: %v", err)
    }
    return nil
}
func IsTokenBlacklisted(token string) (bool, error) {
	// Construct the key that was used to store the blacklisted token
	key := fmt.Sprintf("blacklist:%s", token)

	// Use the Exists command to check if the key exists.
	// The result is the number of keys found (0 or 1 in this case).
	val, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check Redis for blacklisted token: %v", err)
	}

	// If the value is greater than 0, the key exists and the token is blacklisted.
	return val > 0, nil
}
