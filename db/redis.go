package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         "redis:6379",
		Password:     "Sodiq111",
		PoolSize:     50,              // Reduced for local setup
		MinIdleConns: 10,              // Reduced
		PoolTimeout:  3 * time.Second, // Reduced for faster failure
		DialTimeout:  2 * time.Second, // Added for connection speed
		ReadTimeout:  2 * time.Second, // Added for read operations
		WriteTimeout: 2 * time.Second, // Added for write operations
		DB:           0,
	})

	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
