package util

import (
	"crypto/rand"
	"encoding/hex"
)

const base26 = "abcdefghijklmnopqrstuvwxyz"

func Encode(num int) string {
	numStr := ""
	for num > 0 {
		leftover := num % 26
		numStr = string(base26[leftover]) + numStr
		num = num / 26
	}
	return numStr
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
