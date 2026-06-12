package utils

import (
	"math/rand"
	"time"
)

func GenerateTelegramCode() string {
	rand.Seed(time.Now().UnixNano())

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 6)
	
	for i := range code {

		code[i] = chars[rand.Intn(len(chars))]

	}

	return string(code)
}