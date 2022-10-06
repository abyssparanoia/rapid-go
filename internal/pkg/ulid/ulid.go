package ulid

import (
	"math/rand"

	"github.com/oklog/ulid"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/now"
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
