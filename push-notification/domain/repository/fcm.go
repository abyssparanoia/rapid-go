package repository

import "context"

// Fcm ... fcm repository interface
type Fcm interface {
	SubscribeTopic(
		ctx context.Context,
		appID string,
		topic string,
		tokens []string) error
}
