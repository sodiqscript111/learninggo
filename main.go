package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"learninggo/db"
	"learninggo/models"
	"learninggo/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	db.ConnectRedis()

	db.ConnectDatabase()

	server := gin.Default()
	server.DELETE("/event/:id", DeleteEvent)
	server.POST("/register", AddUser)
	server.POST("/login", ValidateUser)
	server.GET("/user/:id", GetUser)
	server.DELETE("/user/:id", DeleteUser)
	server.POST("/createevent", AddEvents)
	server.GET("/events", GetEvents)
	server.GET("/event/:id", GetEvent)
	server.POST("/events/:id/register", AddUserForEvent)
	server.GET("/events/:id/registrations", GetAllRegisters)
	server.GET("/events/:id/registration", GetRegistersWithUserId)
	server.POST("/event/:id/link", CreateEventLink)
	server.POST("/events/register/:code", RegisterHandler)
	server.GET("/event/:id/link", GetAllEventLinks)

	server.PUT("/event/:id/link", DeleteEventLink)

	server.Run(":8080")

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

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	// Cap limit to prevent excessive load
	if limit > 100 {
		limit = 100
	}

	events, err := models.GetAllEvents(page, limit)
	if err != nil {
		log.Printf("GetAllEvents error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "public, max-age=180")
	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"page":   page,
		"limit":  limit,
	})
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

type RegistrationInput struct {
	UserID        string `json:"userId" binding:"required,uuid"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
}

func AddUserForEvent(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	var input RegistrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided: " + err.Error()})
		return
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	registration := db.Registration{
		UserID:        userID,
		EventID:       eventID,
		Status:        input.Status,
		PaymentStatus: input.PaymentStatus,
		RegisteredAt:  time.Now(),
		TicketToken:   utils.GenerateTicketToken(),
	}

	err = models.CreateRegistration(registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create registration: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Registration created successfully",
		"registration": registration,
	})
}

func GetAllRegisters(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}
	registers, err := models.GetAllRegistrations(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"registrations": registers,
	})
}

func GetRegistersWithUserId(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}
	registers, err := models.GetRegistrationByUserId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to get all registrations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"registrations": registers,
		"message":       "Registration created successfully",
	})
}

type eventlinkInput struct {
	UserID    string     `json:"userId" binding:"required,uuid"`
	ExpiresAt *time.Time `json:"expires_at"`
	MaxUses   *int       `json:"max_uses" binding:"required"`
}

func CreateEventLink(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
	}

	code := utils.GenerateCode()
	var input eventlinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	link := db.EventLink{
		Code:      code,
		ExpiresAt: input.ExpiresAt,
		MaxUses:   input.MaxUses,
		CreatedBy: userID,
		EventID:   eventID,
	}
	err = models.AddEventLink(link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Event link created successfully",
	})
	return
}
func GetAllEventLinks(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}
	eventLink, err := models.GetEventLink(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Event link created successfully",
		"Events":  eventLink,
	})
}

func DeleteEventLink(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid eventID format"})
	}
	err = models.DeleteEventLink(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Event link deleted successfully",
	})
}

func RegisterHandler(c *gin.Context) {
	code := c.Param("code")

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	err := models.AddRegistration(code, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}
