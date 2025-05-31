package models

import (
	"errors"
	"github.com/google/uuid"
	"learninggo/db"
	"learninggo/utils"
	"time"
)

func AddRegistration(code string, userID uuid.UUID) error {
	// 1. Find the event by code
	var event db.Event
	if err := db.DB.Where("code = ?", code).First(&event).Error; err != nil {
		return errors.New("event not found")
	}

	var existing db.Registration
	err := db.DB.Where("user_id = ? AND event_id = ?", userID, event.ID).First(&existing).Error
	if err == nil {
		return errors.New("user already registered for this event")
	}

	ticketToken := utils.GenerateTicketToken()

	newRegistration := db.Registration{
		ID:            uuid.New(),
		UserID:        userID,
		EventID:       event.ID,
		RegisteredAt:  time.Now(),
		Status:        "confirmed",
		TicketToken:   ticketToken,
		PaymentStatus: "unpaid",
	}

	if err := db.DB.Create(&newRegistration).Error; err != nil {
		return errors.New("failed to register for event")
	}

	return nil
}
