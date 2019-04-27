package firebaseauth

import (
	"context"
	"net/http"
)

// Service ... service inteface for firebase authentication
type Service interface {
	SetCustomClaims(ctx context.Context, userID string, claims Claims) error
	Authentication(ctx context.Context, r *http.Request) (string, Claims, error)
}
