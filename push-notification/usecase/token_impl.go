package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/push-notification/config"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/repository"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/service"
	"github.com/abyssparanoia/rapid-go/push-notification/usecase/input"
)

type register struct {
	fcmRepository   repository.Fcm
	tokenRepository repository.Token
	tokenService    service.Token
}

func (u *register) SetToken(
	ctx context.Context,
	dto *input.TokenSet) error {

	token := model.NewToken(dto.Platform, dto.AppID, dto.DeviceID, dto.Token)
	err := u.tokenService.Set(ctx, token)
	if err != nil {
		log.Errorm(ctx, "s.tokenService.Set", err)
		return nil
	}

	err = u.fcmRepository.SubscribeTopic(ctx, config.TopicAll, []string{token.Value})
	if err != nil {
		log.Errorm(ctx, "s.fcmRepository.SubscribeTopic", err)
		return nil
	}

	return nil
}

func (u *register) DeleteToken(
	ctx context.Context,
	dto *input.TokenDelete) error {
	token, err := u.tokenRepository.GetByPlatformAndDeviceIDAndUserID(ctx, dto.AppID, dto.UserID, dto.DeviceID, dto.Platform)
	if err != nil {
		log.Errorm(ctx, "u.tokenRepository.GetByPlatformAndDeviceIDAndUserID", err)
		return nil
	}
	if token.Exists() {
		err = u.tokenRepository.Delete(ctx, token.ID)
		if err != nil {
			log.Errorm(ctx, "u.tokenRepository.Delete", err)
			return nil
		}

		err = u.fcmRepository.Unsubscribe(ctx, config.TopicAll, []string{token.Value})
		if err != nil {
			log.Errorm(ctx, "u.fcmRepository.Unsubscribe", err)
			return nil
		}
	}
	return nil
}

// NewRegister ... new register usecase
func NewRegister(
	fcmRepository repository.Fcm,
	tokenRepository repository.Token,
	tokenService service.Token,
) Register {
	return &register{fcmRepository, tokenRepository, tokenService}
}
