package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"
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
		message *model.Message) error
	SendMessageByTopic(
		ctx context.Context,
		appID string,
		topic string,
		message *model.Message) error
}
