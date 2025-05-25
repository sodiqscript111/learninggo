package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"learninggo/db"
	"learninggo/models"
)

func main() {
	db.ConnectDatabase()

	server := gin.Default()

	server.POST("/addUser", func(c *gin.Context) {
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

		// On success, return 201 Created with user info (mask password in production)
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	})

	// Start server on port 8080
	server.Run(":8080")
}
