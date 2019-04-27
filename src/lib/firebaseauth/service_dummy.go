package firebaseauth

import (
	"context"
	"net/http"
)

type dummyService struct {
}

// SetCustomClaims ... set custom claims
func (s *dummyService) SetCustomClaims(ctx context.Context, userID string, claims Claims) error {
	return nil
}

// Authentication ... authenticate
func (s *dummyService) Authentication(ctx context.Context, r *http.Request) (string, Claims, error) {
	userID := "DUMMY_USER_ID"
	claims := Claims{}
	return userID, claims, nil
}

// NewDummyService ... get dummy service
func NewDummyService() Service {
	return &dummyService{}
}
