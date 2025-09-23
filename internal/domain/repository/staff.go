package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/volatiletech/null/v8"
)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Staff interface {
	Get(
		ctx context.Context,
		query GetStaffQuery,
	) (*model.Staff, error)
	List(
		ctx context.Context,
		query ListStaffQuery,
	) (model.Staffs, error)
	Count(
		ctx context.Context,
		query ListStaffQuery,
	) (uint64, error)
	Create(
		ctx context.Context,
		staff *model.Staff,
	) error
	BatchCreate(
		ctx context.Context,
		staffs model.Staffs,
	) error
	Update(
		ctx context.Context,
		staff *model.Staff,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}

type GetStaffQuery struct {
	BaseGetOptions
	ID      null.String
	AuthUID null.String
}

type ListStaffQuery struct {
	BaseListOptions
	TenantID null.String
}
