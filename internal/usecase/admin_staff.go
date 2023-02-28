package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type AdminStaffInteractor interface {
	Create(
		ctx context.Context,
		param *input.AdminCreateStaff,
	) (*model.Staff, error)
}
