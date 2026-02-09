package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}
