package main

import (
	"github.com/gin-gonic/gin"

	"learninggo/db"
	"learninggo/models"
	"net/http"
)

func main() {
	db.InitDB()
	server := gin.Default()

	// Register the GET route for "/home"
	server.GET("/events", getEvents)
	server.POST("events", createEvent)
	// Start the server on port 8080
	server.Run(":8080")
}

func getEvents(context *gin.Context) {
	events := models.GetAllEvents()
	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not bind JSON to struct"})
		return
	}

	event.ID = 1
	event.UserID = 1
	event.Save()
	context.JSON(http.StatusCreated, gin.H{"message": "Event created successfully"})
}
