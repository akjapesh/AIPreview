package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// Initialize Redis connection
func ConnectToRedis() error {
	redisAddr := os.Getenv("REDIS_PORT") // Redis container hostname
	if redisAddr == "" {
		redisAddr = "redis:6379" // Default for Docker
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // Set if Redis requires a password
		DB:       0,  // Use default DB
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}
