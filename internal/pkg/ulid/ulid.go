package ulid

import (
	"math/rand"

	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/oklog/ulid"
)

var New = func() string {
	now := now.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(now), entropy).String()
}

func Mock() string {
	mockULID := "mock"
	New = func() string {
		return mockULID
	}
	return mockULID
}
