package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"learninggo/db"
	"time"
)

var Ctx = context.Background()

func AddEventLink(eventlink db.EventLink) error {

	if err := db.DB.Create(&eventlink).Error; err != nil {
		fmt.Println("DB error:", err)
		return err
	}

	jsonBytes, err := json.Marshal(eventlink)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return err
	}

	redisKey := fmt.Sprintf("eventlink:%v", eventlink.ID)

	if err := db.RedisClient.Set(Ctx, redisKey, jsonBytes, 10*time.Minute).Err(); err != nil {
		fmt.Println("Redis SET error:", err)
		return err
	}

	return nil
}

func GetEventLink(eventlinkId uuid.UUID) ([]db.EventLink, error) {
	var eventlink db.EventLink
	redisKey := fmt.Sprintf("eventlink:%v", eventlinkId)

	cached, err := db.RedisClient.Get(Ctx, redisKey).Result()
	if err == nil {

		if err := json.Unmarshal([]byte(cached), &eventlink); err == nil {
			return []db.EventLink{eventlink}, nil
		}
		fmt.Println("Redis unmarshal error:", err)
	} else if err != redis.Nil {

		fmt.Println("Redis GET error:", err)
	}

	result := db.DB.First(&eventlink, "id = ?", eventlinkId)
	if result.Error != nil {
		return []db.EventLink{}, result.Error
	}

	jsonBytes, err := json.Marshal(eventlink)
	if err == nil {
		_ = db.RedisClient.Set(Ctx, redisKey, jsonBytes, 10*time.Minute).Err()
	}

	return []db.EventLink{eventlink}, nil
}

func DeleteEventLink(eventlinkId uuid.UUID) error {
	redisKey := fmt.Sprintf("eventlink:%v", eventlinkId)
	result := db.DB.Delete(&db.EventLink{}, "id = ?", eventlinkId)
	if result.Error != nil {
		return result.Error
	}
	err := db.RedisClient.Del(Ctx, redisKey).Err()
	if err != nil {
		return err
	}
	return nil
}
