package firebaseauth

import (
	"context"
)

// Firebaseauth ... service inteface for firebase authentication
type Firebaseauth interface {
	SetCustomClaims(ctx context.Context, userID string, claims *Claims) error
	Authentication(ctx context.Context, ah string) (string, *Claims, error)
}
