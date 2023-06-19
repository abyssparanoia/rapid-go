package id

import (
	"github.com/lucsky/cuid"
)

var New = func() string {
	return cuid.New()
}

func Mock() string {
	mockID := "mock"
	New = func() string {
		return mockID
	}
	return mockID
}
