package firebaseauth

import "context"

type contextKey string

const userIDContextKey contextKey = "firebaseauth:user_id"

const claimsContextKey contextKey = "firebaseauth:claims"

// GetUserID ... get FirebaseAuthUID from context
func GetUserID(ctx context.Context) string {
	return ctx.Value(userIDContextKey).(string)
}

// GetClaims ... get FirebaseAuthJWTClaims from context
func GetClaims(ctx context.Context) Claims {
	return ctx.Value(claimsContextKey).(Claims)
}

func setUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func setClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}
