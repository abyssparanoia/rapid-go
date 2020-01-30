package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/push-notification/usecase/input"
)

// Message ... message usecase interface
type Message interface {
	SendToUser(ctx context.Context,
		dto *input.MessageSendToUser) error
	SendToMultiUser(ctx context.Context,
		dto *input.MessageSendToMultiUser) error
	SendToAllUser(ctx context.Context,
		dto *input.MessageSendToAllUser) error
}
