package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

const secretKey = "supersecret"

func GenerateToken(email string, userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}
func VerifyToken(tokenString string) (int, error) {
	// Check if token has "Bearer " prefix
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return 0, errors.New("Token format is invalid")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println("Token parsing error:", err)
		return 0, errors.New("Invalid token")
	}

	if !parsedToken.Valid {
		return 0, errors.New("Invalid token")
	}

	// Extract the user ID from the token claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	userId, ok := claims["userId"].(float64) // JWT stores numbers as float64
	if !ok {
		return 0, errors.New("User ID not found in token")
	}

	return int(userId), nil
}
