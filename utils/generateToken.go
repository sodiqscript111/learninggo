package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
)

func GenerateTicketToken() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		return uuid.New().String()
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
