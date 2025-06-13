package random

import (
	"math/rand"
)

func Number(length int) string {
	digits := make([]byte, length)
	for i := 0; i < length; i++ {
		digits[i] = byte(rand.Intn(10)) + '0'
	}
	return string(digits)
}

func String(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}
