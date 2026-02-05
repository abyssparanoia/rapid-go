package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type StaffStaffInteractor interface {
	Get(
		ctx context.Context,
		param *input.StaffGetStaff,
	) (*model.Staff, error)
	List(
		ctx context.Context,
		param *input.StaffListStaffs,
	) (*output.ListStaffs, error)
	Create(
		ctx context.Context,
		param *input.StaffCreateStaff,
	) (*model.Staff, error)
	Update(
		ctx context.Context,
		param *input.StaffUpdateStaff,
	) (*model.Staff, error)
}
