package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/repository"
)

type token struct {
	tokenRepository repository.Token
}

func (s *token) Set(ctx context.Context,
	token *model.Token) error {

	currentToken, err := s.tokenRepository.GetByPlatformAndDeviceIDAndUserID(ctx, token.AppID, token.UserID, token.DeviceID, token.Platform)
	if err != nil {
		log.Errorm(ctx, "s.tokenRepository.GetByPlatformAndDeviceIDAndUserID", err)
		return err
	}
	// if current token exists, check token value and update token
	if currentToken.Exists() && token.Value != currentToken.Value {
		token.Value = currentToken.Value
		err = s.tokenRepository.Update(ctx, token)
		if err != nil {
			log.Errorm(ctx, "s.tokenRepository.Update", err)
			return err
		}
		// if not exists, create token document
	} else if !currentToken.Exists() {
		_, err = s.tokenRepository.Create(ctx, token)
		if err != nil {
			log.Errorm(ctx, "s.tokenRepository.Create", err)
			return err
		}
	}

	return nil
}

// NewToken ... new token service
func NewToken(
	tokenRepository repository.Token,
) Token {
	return &token{tokenRepository}
}
