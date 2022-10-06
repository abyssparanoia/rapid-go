package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type UserInteractor interface {
	CreateRoot(
		ctx context.Context,
		param *input.CreateRootUser,
	) error
}
