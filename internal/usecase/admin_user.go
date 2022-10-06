package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type AdminUserInteractor interface {
	Create(
		ctx context.Context,
		param *input.AdminCreateUser,
	) (*model.User, error)
}
