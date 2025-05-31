package models

import (
	"fmt"
	"github.com/google/uuid"
	"learninggo/db"
)

func AddEventLink(eventlink db.EventLink) error {
	err := db.DB.Create(&eventlink)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func GetEventLink(eventlinkId uuid.UUID) ([]db.EventLink, error) {
	var eventlinks db.EventLink
	result := db.DB.First(&eventlinks, "id = ?", eventlinkId)
	if result.Error != nil {
		return []db.EventLink{}, result.Error
	}
	return []db.EventLink{eventlinks}, nil
}

func DeleteEventLink(eventlinkId uuid.UUID) error {
	result := db.DB.Delete(&db.EventLink{}, "id = ?", eventlinkId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
