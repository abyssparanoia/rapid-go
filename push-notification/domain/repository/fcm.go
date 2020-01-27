package repository

import "context"

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
}
