package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"learninggo/db"
)

func AddUser(user db.User) error {
	result := db.DB.Create(&user)
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
	result := db.DB.First(&user, "id = ?", userID)

	if result.Error != nil {
		return db.User{}, result.Error
	}
	return user, nil
}

func DeleteUser(userId uuid.UUID) error {
	result := db.DB.Delete(&db.User{}, "id = ?", userId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
