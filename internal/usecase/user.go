package usecase

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
)

type UserInteractor interface {
	CreateRoot(
		ctx context.Context,
		param *input.CreateRootUser,
	) error
}
