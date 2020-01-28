package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/push-notification/usecase/input"
)

// Register ... register usecase
type Register interface {
	SetToken(
		ctx context.Context,
		dto *input.RegisterSetToken) error
}
