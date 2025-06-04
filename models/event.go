package models

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"time"

	"github.com/google/uuid"
	"learninggo/db"
)

func AddEvent(event db.Event) error {
	redisKey := fmt.Sprintf("event:%s", event.ID.String())
	result := db.DB.Create(&event)
	if result.Error != nil {
		return result.Error
	}
	jsonBytes, err := json.Marshal(&event)
	if err != nil {
		return err
	}
	err = db.RedisClient.Set(Ctx, redisKey, jsonBytes, 2*time.Minute).Err()
	if err != nil {
		fmt.Println("Redis set error:", err)
		return err
	}
	// Invalidate all events cache
	db.RedisClient.Del(Ctx, "events:all")
	return nil
}

func GetAllEvents(page, limit int) ([]db.Event, error) {
	var events []db.Event
	redisKey := fmt.Sprintf("events:page:%d:limit:%d", page, limit)

	// Try to fetch from Redis
	cached, err := db.RedisClient.Get(Ctx, redisKey).Bytes()
	if err == nil {
		// Decompress cached data
		reader, err := gzip.NewReader(bytes.NewReader(cached))
		if err == nil {
			var decompressed []byte
			decompressed, err = io.ReadAll(reader)
			reader.Close()
			if err == nil && json.Unmarshal(decompressed, &events) == nil {
				return events, nil
			}
		}
	} else if err != redis.Nil {
		fmt.Println("Redis get error:", err)
	}

	// Fetch from database if cache miss
	offset := (page - 1) * limit
	result := db.DB.Select("id", "title", "start_time", "end_time", "location", "max_capacity", "is_paid", "price", "created_by").
		Order("start_time DESC"). // Added for consistency and index usage
		Offset(offset).Limit(limit).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}

	// Cache the result with compression if Redis is available
	data, err := json.Marshal(events)
	if err == nil {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err = writer.Write(data)
		writer.Close()
		if err == nil {
			err = db.RedisClient.Set(Ctx, redisKey, buf.Bytes(), 3*time.Minute).Err() // Reduced TTL
			if err != nil && err != redis.ErrClosed {
				fmt.Println("Redis set error:", err)
			}
		}
	}

	return events, nil
}

func GetEvent(eventId uuid.UUID) (db.Event, error) {
	var event db.Event
	redisKey := fmt.Sprintf("event:%s", eventId.String())

	// Try to fetch from Redis
	cached, err := db.RedisClient.Get(Ctx, redisKey).Bytes()
	if err == nil {
		if err := json.Unmarshal(cached, &event); err == nil {
			return event, nil
		}
	} else if err != redis.Nil {
		fmt.Println("Redis get error:", err)
	}

	// Fetch from database if cache miss
	result := db.DB.Select("id", "title", "start_time", "end_time", "location", "max_capacity", "is_paid", "price", "created_by").
		First(&event, "id = ?", eventId)
	if result.Error != nil {
		return db.Event{}, result.Error
	}

	// Cache the result if Redis is available
	data, err := json.Marshal(event)
	if err == nil {
		err = db.RedisClient.Set(Ctx, redisKey, data, 10*time.Minute).Err()
		if err != nil && err != redis.ErrClosed {
			fmt.Println("Redis set error:", err)
		}
	}

	return event, nil
}

func DeleteEvent(eventId uuid.UUID) error {
	redisKey := fmt.Sprintf("event:%s", eventId.String())
	result := db.DB.Delete(&db.Event{}, "id = ?", eventId)
	if result.Error != nil {
		return result.Error
	}
	// Delete event from Redis
	err := db.RedisClient.Del(Ctx, redisKey).Err()
	if err != nil && err != redis.ErrClosed {
		fmt.Println("Redis delete error:", err)
	}
	// Invalidate all events cache
	err = db.RedisClient.Del(Ctx, "events:all").Err()
	if err != nil && err != redis.ErrClosed {
		fmt.Println("Redis delete error for all events:", err)
	}
	return nil
}
