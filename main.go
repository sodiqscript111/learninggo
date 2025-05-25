package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"learninggo/db"
	"learninggo/models"
	"net/http"
)

func main() {
	db.ConnectDatabase()

	server := gin.Default()

	server.POST("/register", AddUser)
	server.POST("/login", ValidateUser)
	server.GET("/user/:id", GetUser)
	server.DELETE("/user/:id", DeleteUser)
	server.POST("/createevent", AddEvents)
	server.GET("/events", GetEvents)
	server.GET("/event/:id", GetEvent)
	server.Run(":8080")
	server.DELETE("/event/:id", DeleteEvent)
}

func AddUser(c *gin.Context) {
	var user db.User

	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON provided: " + err.Error(),
		})
		return
	}

	if err := models.AddUser(user); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func ValidateUser(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON provided: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
func GetUser(c *gin.Context) {
	userIDParam := c.Param("id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := models.GetUser(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "User found",
		"user":    user,
	})
}

func DeleteUser(c *gin.Context) {
	userIDParam := c.Param("id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	if err := models.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User deleted successfully",
	})
}

func AddEvents(c *gin.Context) {
	var event db.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided: " + err.Error()})
		return
	}

	var user db.User
	if err := db.DB.First(&user, "id = ?", event.CreatedBy).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with ID " + event.CreatedBy.String() + " does not exist"})
		return
	}

	if err := models.AddEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

func GetEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func GetEvent(c *gin.Context) {
	eventIDParam := c.Param("id")

	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := models.GetEvent(eventID)
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event": event,
	})
}

func DeleteEvent(c *gin.Context) {
	eventIDParam := c.Param("id")

	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	err = models.DeleteEvent(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Event deleted successfully",
	})
}
