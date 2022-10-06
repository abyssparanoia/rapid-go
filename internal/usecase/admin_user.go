package usecase

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
)

type AdminUserInteractor interface {
	Create(
		ctx context.Context,
		param *input.AdminCreateUser,
	) (*model.User, error)
}
