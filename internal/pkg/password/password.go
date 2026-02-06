package password

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const DefaultLength = 16

// Generate is a variable function for password generation (mockable)
var Generate = func(length int) (string, error) {
	byteLength := max((length*3)/4, 1)

	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}

// Mock replaces Generate with a function that returns a fixed value
func Mock(mockPassword string) string {
	Generate = func(_ int) (string, error) {
		return mockPassword, nil
	}
	return mockPassword
}
