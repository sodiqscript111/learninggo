package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sony/gobreaker"
	_ "gorm.io/gorm"
	"learninggo/db"
	"log"
	"sync"
	"time"
)

// Circuit breaker for database operations
var (
	cb     *gobreaker.CircuitBreaker
	cbOnce sync.Once
)

func getCircuitBreaker() *gobreaker.CircuitBreaker {
	cbOnce.Do(func() {
		cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "db-create-registration",
			MaxRequests: 10,
			Interval:    30 * time.Second,
			Timeout:     5 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 5
			},
		})
	})
	return cb
}

// CreateRegistrations handles batch creation of registrations with async Redis caching
func CreateRegistrations(registrations []db.Registration) error {
	if len(registrations) == 0 {
		return nil
	}

	// Execute within circuit breaker
	_, err := getCircuitBreaker().Execute(func() (interface{}, error) {
		// Start a transaction
		tx := db.DB.Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}

		// Batch insert (100 records per batch)
		if err := tx.CreateInBatches(registrations, 100).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("batch create failed: %w", err)
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("commit failed: %w", err)
		}

		return nil, nil
	})
	if err != nil {
		log.Printf("CreateRegistrations failed: %v", err)
		return err
	}

	// Cache in Redis asynchronously
	go func() {
		for _, reg := range registrations {
			redisKey := fmt.Sprintf("registration:user:%s", reg.UserID)
			// Cache minimal data to reduce serialization overhead
			cacheData := map[string]interface{}{
				"userID":        reg.UserID,
				"eventID":       reg.EventID,
				"status":        reg.Status,
				"paymentStatus": reg.PaymentStatus,
			}
			if jsonBytes, err := json.Marshal(cacheData); err == nil {
				if err := db.RedisClient.Set(context.Background(), redisKey, jsonBytes, 20*time.Minute).Err(); err != nil {
					log.Printf("Failed to cache registration for user %s: %v", reg.UserID, err)
				}
			}
		}
	}()

	return nil
}

// CreateRegistration wraps CreateRegistrations for single record
func CreateRegistration(register db.Registration) error {
	return CreateRegistrations([]db.Registration{register})
}

// GetAllRegistrations retrieves registrations by eventID with caching
func GetAllRegistrations(eventID uuid.UUID) ([]db.Registration, error) {
	redisKey := fmt.Sprintf("registrations:event:%s", eventID)
	var registrations []db.Registration

	// Try Redis cache
	if cached, err := db.RedisClient.Get(context.Background(), redisKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &registrations); err == nil {
			return registrations, nil
		}
	}

	// Query database with preloading
	if err := db.DB.Preload("User").Preload("Event").Where("event_id = ?", eventID).Find(&registrations).Error; err != nil {
		log.Printf("GetAllRegistrations failed for event %s: %v", eventID, err)
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Cache result asynchronously
	if jsonBytes, err := json.Marshal(registrations); err == nil {
		go func() {
			if err := db.RedisClient.Set(context.Background(), redisKey, jsonBytes, 20*time.Minute).Err(); err != nil {
				log.Printf("Failed to cache registrations for event %s: %v", eventID, err)
			}
		}()
	}

	return registrations, nil
}

// GetAllRegistrationsByUserId retrieves registrations by userID with caching
func GetAllRegistrationsByUserId(userID uuid.UUID) ([]db.Registration, error) {
	redisKey := fmt.Sprintf("registrations:user:%s", userID)
	var registrations []db.Registration

	// Try Redis cache
	if cached, err := db.RedisClient.Get(context.Background(), redisKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &registrations); err == nil {
			return registrations, nil
		}
	}

	// Query database with preloading
	if err := db.DB.Preload("User").Preload("Event").Where("user_id = ?", userID).Find(&registrations).Error; err != nil {
		log.Printf("GetAllRegistrationsByUserId failed for user %s: %v", userID, err)
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Cache result asynchronously
	if jsonBytes, err := json.Marshal(registrations); err == nil {
		go func() {
			if err := db.RedisClient.Set(context.Background(), redisKey, jsonBytes, 20*time.Minute).Err(); err != nil {
				log.Printf("Failed to cache registrations for user %s: %v", userID, err)
			}
		}()
	}

	return registrations, nil
}

// GetRegistrationByUserId retrieves a single registration by userID
func GetRegistrationByUserId(userID uuid.UUID) (db.Registration, error) {
	redisKey := fmt.Sprintf("registration:user:%s", userID)
	var registration db.Registration

	// Try Redis cache
	if cached, err := db.RedisClient.Get(context.Background(), redisKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &registration); err == nil {
			return registration, nil
		}
	}

	// Query database
	if err := db.DB.Preload("User").First(&registration, "user_id = ?", userID).Error; err != nil {
		log.Printf("GetRegistrationByUserId failed for user %s: %v", userID, err)
		return db.Registration{}, fmt.Errorf("query failed: %w", err)
	}

	// Cache result asynchronously
	if jsonBytes, err := json.Marshal(registration); err == nil {
		go func() {
			if err := db.RedisClient.Set(context.Background(), redisKey, jsonBytes, 20*time.Minute).Err(); err != nil {
				log.Printf("Failed to cache registration for user %s: %v", userID, err)
			}
		}()
	}

	return registration, nil
}

// DeleteRegistrationCache invalidates cache for a user
func DeleteRegistrationCache(userID uuid.UUID) error {
	redisKey := fmt.Sprintf("registration:user:%s", userID)
	if err := db.RedisClient.Del(context.Background(), redisKey).Err(); err != nil {
		log.Printf("Failed to delete cache for user %s: %v", userID, err)
		return fmt.Errorf("cache delete failed: %w", err)
	}
	return nil
}
