package uuid

import (
	"encoding/base64"

	"github.com/google/uuid"
)

var UUID = func() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

var UUIDBase64 = func() string {
	uuid, _ := uuid.NewRandom()
	uuidBinary, _ := uuid.MarshalBinary()
	return base64.RawURLEncoding.EncodeToString(uuidBinary)
}

func MockUUIDBase64() string {
	mockUUIDBase64 := "mock"
	UUIDBase64 = func() string {
		return mockUUIDBase64
	}
	return mockUUIDBase64
}
