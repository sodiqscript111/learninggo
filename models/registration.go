package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"learninggo/db"
	"time"
)

func CreateRegistration(register db.Registration) error {

	result := db.DB.Create(&register)
	if result.Error != nil {
		return result.Error
	}

	redisKey := fmt.Sprintf("registration:user:%s", register.UserID)
	jsonBytes, err := json.Marshal(register)
	if err != nil {
		return err
	}

	err = db.RedisClient.Set(Ctx, redisKey, jsonBytes, 20*time.Minute).Err()
	if err != nil {

		fmt.Println("Failed to cache registration:", err)
	}

	return nil
}

func GetAllRegistrations(eventID uuid.UUID) ([]db.Registration, error) {
	var registrations []db.Registration

	result := db.DB.Where("event_id = ?", eventID).Find(&registrations)
	if result.Error != nil {
		return nil, result.Error
	}
	return registrations, nil
}

func GetAllRegistrationsByUserId(userID uuid.UUID) ([]db.Registration, error) {
	var registrations []db.Registration

	result := db.DB.Where("user_id = ?", userID).Find(&registrations)
	if result.Error != nil {
		return nil, result.Error
	}
	return registrations, nil
}

func GetRegistrationByUserId(userID uuid.UUID) (db.Registration, error) {
	var registration db.Registration
	redisKey := fmt.Sprintf("registration:user:%s", userID)

	cached, err := db.RedisClient.Get(Ctx, redisKey).Bytes()
	if err == nil {
		if err := json.Unmarshal(cached, &registration); err == nil {
			return registration, nil
		}

	}

	result := db.DB.Where("user_id = ?", userID).First(&registration)
	if result.Error != nil {
		return db.Registration{}, result.Error
	}

	jsonBytes, err := json.Marshal(registration)
	if err == nil {
		_ = db.RedisClient.Set(Ctx, redisKey, jsonBytes, 20*time.Minute).Err()
	}

	return registration, nil
}
