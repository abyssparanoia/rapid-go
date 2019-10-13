package firebaseauth

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/src/lib/log"
)

type service struct {
}

// SetCustomClaims ... set custom claims
func (s *service) SetCustomClaims(ctx context.Context, userID string, claims *Claims) error {
	c, err := getAuthClient(ctx)
	if err != nil {
		log.Errorm(ctx, "getAuthClient", err)
		return err
	}

	err = c.SetCustomUserClaims(ctx, userID, claims.ToMap())
	if err != nil {
		log.Errorm(ctx, "c.SetCustomUserClaims", err)
		return err
	}

	return nil
}

// Authentication ... authenticate
func (s *service) Authentication(ctx context.Context, ah string) (string, *Claims, error) {
	var userID string
	claims := &Claims{}

	c, err := getAuthClient(ctx)
	if err != nil {
		log.Warningm(ctx, "getAuthClient", err)
		return userID, claims, err
	}

	token := getTokenByAuthHeader(ah)
	if token == "" {
		err := log.Warninge(ctx, "token empty error")
		return userID, claims, err
	}

	t, err := c.VerifyIDToken(ctx, token)
	if err != nil {
		msg := fmt.Sprintf("c.VerifyIDToken: %s", token)
		log.Warningm(ctx, msg, err)
		return userID, claims, err
	}

	userID = t.UID
	claims.SetMap(t.Claims)

	return userID, claims, nil
}

// NewService ... get service
func NewService() Service {
	return &service{}
}
