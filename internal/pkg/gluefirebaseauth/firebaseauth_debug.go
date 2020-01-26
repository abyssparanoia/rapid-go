package gluefirebaseauth

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
)

type firebaseauthDebug struct {
	cli *auth.Client
}

// Authentication ... authenticate
func (s *firebaseauthDebug) Authentication(ctx context.Context, ah string) (string, *Claims, error) {
	var userID string
	claims := &Claims{}

	// ユーザーを取得できたらデバッグリクエストと判定する
	if user := getUserByAuthHeader(ah); user != "" {
		claims = newDummyClaims()
		return user, claims, nil
	}

	// 通常の認証を行う
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

func (s *firebaseauthDebug) CreateTokenWithClaims(ctx context.Context, userID string, claims *Claims) (string, error) {
	token, err := s.cli.CustomTokenWithClaims(ctx, userID, claims.ToMap())
	if err != nil {
		log.Errorm(ctx, "s.cli.CustomTokenWithClaims", err)
		return "", err
	}
	return token, nil
}

func (s *firebaseauthDebug) CreateUser(ctx context.Context, email string, password string) (*auth.UserRecord, error) {

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

func (s *firebaseauthDebug) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {

	userRecord, err := s.cli.GetUserByEmail(ctx, email)
	if err != nil {
		if userRecord == nil {
			return nil, nil
		}
		return nil, err
	}

	return userRecord, nil

}

// NewDebug ... Debuggluefirebaseauthを作成する
func NewDebug(cli *auth.Client) Firebaseauth {
	return &firebaseauthDebug{cli}
}
