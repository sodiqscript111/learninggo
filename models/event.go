package models

import (
	"github.com/google/uuid"
	"learninggo/db"
)

func AddEvent(event db.Event) error {
	result := db.DB.Create(&event)
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
	result := db.DB.First(&db.Event{}, eventId)
	if result.Error != nil {
		return db.Event{}, result.Error
	}
	return db.Event{}, nil
}

func DeleteEvent(EventId uuid.UUID) error {
	result := db.DB.Delete(&db.User{}, "id = ?", EventId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
