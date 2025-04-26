package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

func VerifyToken(token string) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return errors.New("Invalid token")
	}

	tokenIsVlid := parsedToken.Valid
	if !tokenIsVlid {
		return errors.New("Invalid token")
	}

	//claims, ok := parsedToken.Claims.(jwt.MapClaims)
	//if !ok {
	//return errors.New("Invalid token claims")
	//}

	//email := claims["email"].(string)
	//userId := claims["userId"].(int)
	return nil
}
