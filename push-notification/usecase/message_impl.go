package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/repository"
	"github.com/abyssparanoia/rapid-go/push-notification/usecase/input"
)

type message struct {
	fcmRepository   repository.Fcm
	tokenRepository repository.Token
}

func (u *message) SendToUser(ctx context.Context,
	dto *input.MessageSendToUser) error {

	tokens, err := u.tokenRepository.ListByUserID(ctx, dto.AppID, dto.UserID)
	if err != nil {
		log.Errorm(ctx, "u.tokenRepository.ListByUserID", err)
		return nil
	}
	if len(tokens) == 0 {
		log.Warningf(ctx, "no regist tokens user: %s", dto.UserID)
		return nil
	}

	tokenValues := make([]string, len(tokens))
	for index, token := range tokens {
		tokenValues[index] = token.Value
	}

	err = u.fcmRepository.SendMessageByTokens(ctx, dto.AppID, tokenValues, dto.Message)
	if err != nil {
		log.Errorm(ctx, "u.fcmRepository.SendMessageByTokens", err)
		return nil
	}

	return nil
}

// NewMessage ... new message usecase
func NewMessage(fcmRepository repository.Fcm, tokenRepository repository.Token) Message {
	return &message{fcmRepository, tokenRepository}
}
