package firebaseauth

import (
	"context"
)

// Service ... service inteface for firebase authentication
type Service interface {
	SetCustomClaims(ctx context.Context, userID string, claims *Claims) error
	Authentication(ctx context.Context, ah string) (string, *Claims, error)
}
