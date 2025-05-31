package models

import (
	"github.com/google/uuid"
	"learninggo/db"
)

func CreateRegistration(register db.Registration) error {
	result := db.DB.Create(&register)
	if result.Error != nil {
		return result.Error
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

func GetAllRegistrationWithUserId(userID uuid.UUID) ([]db.Registration, error) {
	var registrations []db.Registration
	result := db.DB.Where("user_id = ?", userID).Find(&registrations)
	if result.Error != nil {
		return nil, result.Error
	}
	return registrations, nil
}
