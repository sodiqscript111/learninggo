package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Sodiq111",
		DB:       0,
	})

	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
