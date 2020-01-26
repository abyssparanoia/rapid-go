package gluefirebaseauth

import (
	"context"

	"firebase.google.com/go/auth"
)

// Firebaseauth ... service inteface for firebase authentication
type Firebaseauth interface {
	CreateTokenWithClaims(ctx context.Context, userID string, claims *Claims) (string, error)
	Authentication(ctx context.Context, ah string) (string, *Claims, error)
	GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
	CreateUser(ctx context.Context, email string, password string) (*auth.UserRecord, error)
}
