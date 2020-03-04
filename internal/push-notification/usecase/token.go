package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/push-notification/usecase/input"
)

// Token ... token usecase
type Token interface {
	Set(ctx context.Context,
		dto *input.TokenSet) error
	Delete(ctx context.Context,
		dto *input.TokenDelete) error
}
