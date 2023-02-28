package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type StaffInteractor interface {
	CreateRoot(
		ctx context.Context,
		param *input.CreateRootStaff,
	) error
}
