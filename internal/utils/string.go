package utils

import (
	"math/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwsyz0123456789"

func RandString(n int) string {
	var sb strings.Builder
	k := len(charset)

	for i := 0; i < n; i++ {
		c := charset[rand.Intn(k)] //nolint:gosec // for faster performance
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandKey generates a random API key only for testing Database
func RandKey() string {
	return RandString(32)
}
