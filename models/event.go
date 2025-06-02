package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"learninggo/db"
	"time"
)

func AddEvent(event db.Event) error {
	redisKey := fmt.Sprintf("event:%v", event.ID)
	result := db.DB.Create(&event)
	jsonBytes, _ := json.Marshal(&event)
	err := db.RedisClient.Set(Ctx, redisKey, jsonBytes, 20*time.Minute).Err()
	if err != nil {
		return err
	}
	return result.Error
}

func GetAllEvents() ([]db.Event, error) {
	var events []db.Event

	result := db.DB.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func GetEvent(eventId uuid.UUID) (db.Event, error) {
	var event db.Event
	redisKey := fmt.Sprintf("event:%s", eventId.String())

	cached, err := db.RedisClient.Get(Ctx, redisKey).Bytes()
	if err == nil {

		if err := json.Unmarshal(cached, &event); err == nil {
			return event, nil
		}
	}

	result := db.DB.First(&event, "id = ?", eventId)
	if result.Error != nil {
		return db.Event{}, result.Error
	}

	data, err := json.Marshal(event)
	if err == nil {
		_ = db.RedisClient.Set(Ctx, redisKey, data, 10*time.Minute).Err()
	}

	return event, nil
}

func DeleteEvent(EventId uuid.UUID) error {
	redisKey := fmt.Sprintf("event:%s", EventId.String())
	result := db.DB.Delete(&db.User{}, "id = ?", EventId)
	if result.Error != nil {
		return result.Error
	}
	err := db.RedisClient.Del(Ctx, redisKey).Err()
	if err != nil {
		fmt.Println("Redis del error:", err)
		return err
	}

	return nil
}
