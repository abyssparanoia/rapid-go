package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/push-notification/usecase/input"
)

// Register ... register usecase
type Register interface {
	SetToken(
		ctx context.Context,
		dto *input.TokenSet) error
	DeleteToken(
		ctx context.Context,
		dto *input.TokenDelete) error
}
