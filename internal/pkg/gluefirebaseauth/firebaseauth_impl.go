package gluefirebaseauth

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"

	"firebase.google.com/go/auth"
)

type firebaseauth struct {
	cli *auth.Client
}

func (s *firebaseauth) CreateTokenWithClaims(ctx context.Context, userID string, claims *Claims) (string, error) {
	token, err := s.cli.CustomTokenWithClaims(ctx, userID, claims.ToMap())
	if err != nil {
		log.Errorm(ctx, "s.cli.CustomTokenWithClaims", err)
		return "", err
	}
	return token, nil
}

// Authentication ... authenticate
func (s *firebaseauth) Authentication(ctx context.Context, ah string) (string, *Claims, error) {
	var userID string
	claims := &Claims{}

	token := getTokenByAuthHeader(ah)
	if token == "" {
		err := log.Warninge(ctx, "token empty error")
		return userID, claims, err
	}

	t, err := s.cli.VerifyIDToken(ctx, token)
	if err != nil {
		msg := fmt.Sprintf("c.VerifyIDToken: %s", token)
		log.Warningm(ctx, msg, err)
		return userID, claims, err
	}

	userID = t.UID
	claims.SetMap(t.Claims)

	return userID, claims, nil
}

func (s *firebaseauth) CreateUser(ctx context.Context, email string, password string) (*auth.UserRecord, error) {

	userCreate := &auth.UserToCreate{}
	userCreate = userCreate.Email(email)
	userCreate = userCreate.Password(password)

	userRecord, err := s.cli.CreateUser(ctx, userCreate)
	if err != nil {
		log.Errorm(ctx, "s.cli.CreateUser", err)
		return nil, err
	}

	return userRecord, nil
}

func (s *firebaseauth) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {

	userRecord, err := s.cli.GetUserByEmail(ctx, email)
	if err != nil {
		if userRecord == nil {
			return nil, nil
		}
		return nil, err
	}

	return userRecord, nil

}

// New ... get firebaseauth
func New(cli *auth.Client) Firebaseauth {
	return &firebaseauth{cli}
}
