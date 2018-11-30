package firebaseauth

import (
	"context"
	"net/http"
)

// Service ... Firebase認証の機能を提供する
type Service interface {
	SetCustomClaims(ctx context.Context, userID string, claims Claims) error
	Authentication(ctx context.Context, r *http.Request) (string, Claims, error)
}
