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
	fcmRepository repository.Fcm
	tokenService  service.Token
}

func (s *register) SetToken(
	ctx context.Context,
	dto *input.RegisterSetToken) error {

	token := model.NewToken(dto.Platform, dto.AppID, dto.DeviceID, dto.Token)
	err := s.tokenService.Set(ctx, token)
	if err != nil {
		log.Errorm(ctx, "s.tokenService.Set", err)
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
	fcmRepository repository.Fcm,
	tokenService service.Token,
) Register {
	return &register{fcmRepository, tokenService}
}
