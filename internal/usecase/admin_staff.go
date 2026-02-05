package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type AdminStaffInteractor interface {
	Get(
		ctx context.Context,
		param *input.AdminGetStaff,
	) (*model.Staff, error)
	List(
		ctx context.Context,
		param *input.AdminListStaffs,
	) (*output.ListStaffs, error)
	Create(
		ctx context.Context,
		param *input.AdminCreateStaff,
	) (*model.Staff, error)
	Update(
		ctx context.Context,
		param *input.AdminUpdateStaff,
	) (*model.Staff, error)
}
