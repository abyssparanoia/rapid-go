package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
)

// Fcm ... fcm repository interface
type Fcm interface {
	SubscribeTopic(
		ctx context.Context,
		topic string,
		tokens []string) error
	Unsubscribe(
		ctx context.Context,
		topic string,
		tokens []string) error
	SendMessageByTokens(
		ctx context.Context,
		appID string,
		tokens []string,
		src *model.Message) error
}
