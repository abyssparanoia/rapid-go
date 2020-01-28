package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/push-notification/config"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/repository"
	"github.com/abyssparanoia/rapid-go/push-notification/usecase/input"
)

type register struct {
	tokenRepository repository.Token
	fcmRepository   repository.Fcm
}

func (s *register) SetToken(
	ctx context.Context,
	dto *input.RegisterSetToken) error {

	token := model.NewToken(dto.Platform, dto.AppID, dto.DeviceID, dto.Token)
	_, err := s.tokenRepository.Create(ctx, token)
	if err != nil {
		log.Errorm(ctx, "s.tokenRepository.Create", err)
		return nil
	}

	err = s.fcmRepository.SubscribeTopic(ctx, config.TopicAll, []string{token.Value})
	if err != nil {
		log.Errorm(ctx, "s.fcmRepository.SubscribeTopic", err)
		return nil
	}

	return nil
}

// NewRegister ... new register usecase
func NewRegister(
	tokenRepository repository.Token,
	fcmRepository repository.Fcm,
) Register {
	return &register{tokenRepository, fcmRepository}
}
