package Middlewares

import (
	"github.com/gin-gonic/gin"
	"learninggo/utils"
	"net/http"
)

func Authorize(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	// Extract user ID from the token
	userId, err := utils.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Set("userId", userId)
	c.Next()
}
