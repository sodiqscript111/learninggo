package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateCode() string {
	b := make([]byte, 16)
	rand.Read(b)
	fmt.Printf("%x\n", b)
	return hex.EncodeToString(b)

}
