package helper

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
