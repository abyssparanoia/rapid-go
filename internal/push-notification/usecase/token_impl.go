package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/push-notification/config"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/usecase/input"
)

type token struct {
	fcmRepository   repository.Fcm
	tokenRepository repository.Token
	tokenService    service.Token
}

func (u *token) Set(ctx context.Context,
	dto *input.TokenSet) error {

	token := model.NewToken(dto.Platform, dto.AppID, dto.DeviceID, dto.Token)
	err := u.tokenService.Set(ctx, token)
	if err != nil {
		return nil
	}

	err = u.fcmRepository.SubscribeTopic(ctx, config.TopicAll, []string{token.Value})
	if err != nil {
		return nil
	}

	return nil
}

func (u *token) Delete(ctx context.Context,
	dto *input.TokenDelete) error {
	token, err := u.tokenRepository.GetByPlatformAndDeviceIDAndUserID(ctx, dto.AppID, dto.UserID, dto.DeviceID, dto.Platform)
	if err != nil {
		return nil
	}
	if token.Exists() {
		err = u.tokenRepository.Delete(ctx, token.ID)
		if err != nil {
			return nil
		}

		err = u.fcmRepository.Unsubscribe(ctx, config.TopicAll, []string{token.Value})
		if err != nil {
			return nil
		}
	}
	return nil
}

// NewToken ... new token usecase
func NewToken(fcmRepository repository.Fcm,
	tokenRepository repository.Token,
	tokenService service.Token) Token {
	return &token{fcmRepository, tokenRepository, tokenService}
}
