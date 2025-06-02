package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"learninggo/db"
	"time"
)

func AddUser(user db.User) error {
	redisKey := fmt.Sprintf("user:%v", user.ID)
	result := db.DB.Create(&user)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = db.RedisClient.Set(Ctx, redisKey, jsonBytes, time.Hour).Err()
	if err != nil {
		return err
	}
	return result.Error
}
func ValidateUser(user db.User) (bool, error) {
	var existingUser db.User
	result := db.DB.Where("email = ?", user.Email).First(&existingUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {

			return true, nil
		}

		return false, result.Error
	}

	return false, nil
}

func GetUser(userID uuid.UUID) (db.User, error) {
	var user db.User
	redisKey := fmt.Sprintf("user:%s", userID.String())

	cached, err := db.RedisClient.Get(Ctx, redisKey).Bytes()
	if err == nil {
		if err := json.Unmarshal(cached, &user); err == nil {
			return user, nil
		}

	}

	result := db.DB.First(&user, "id = ?", userID)
	if result.Error != nil {
		return db.User{}, result.Error
	}

	data, err := json.Marshal(user)
	if err == nil {
		_ = db.RedisClient.Set(Ctx, redisKey, data, 10*time.Minute).Err()
	}

	return user, nil
}

func DeleteUser(userId uuid.UUID) error {
	redisKey := fmt.Sprintf("user:%s", userId.String())
	err := db.RedisClient.Del(Ctx, redisKey).Err()
	if err != nil {
		return err
	}
	result := db.DB.Delete(&db.User{}, "id = ?", userId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
