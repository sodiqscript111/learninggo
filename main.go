package main

import (
	"github.com/gin-gonic/gin"
	"learninggo/db"
	"learninggo/models"
	"net/http"
	"strconv"
)

func main() {
	db.InitDB()

	server := gin.Default()
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)
	server.PUT("/events/:id", editEvent)
	server.DELETE("/events/:id", deleteEvent)
	server.POST("/events", createEvent)
	server.POST("/signup", createUser)
	server.POST("/login", login)

	server.Run(":8080")
}

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}
func getEvent(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve eventID"})
		return
	}
	event, err := models.GetEventById(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
	}
	c.JSON(http.StatusOK, event)
}

func deleteEvent(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	err = models.DeleteById(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func createEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Don't hardcode ID — database handles it
	// you can later make this dynamic if you have users

	if err := event.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}
func editEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = models.EditById(eventId, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

func createUser(c *gin.Context) {
	var user models.User

	// Parse and validate request body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Insert user into DB
	user.InsertUser()

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.Id,
			"email": user.Email,
		},
	})
}

func login(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	}
	err = user.ValidateUser()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
}
