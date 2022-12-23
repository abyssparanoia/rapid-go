package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Tenant interface {
	Get(
		ctx context.Context,
		id string,
		orFail bool,
	) (*model.Tenant, error)
	List(
		ctx context.Context,
		query ListTenantsQuery,
	) ([]*model.Tenant, error)
	Count(
		ctx context.Context,
		query ListTenantsQuery,
	) (uint64, error)
	Create(
		ctx context.Context,
		tenant *model.Tenant,
	) (*model.Tenant, error)
	Update(
		ctx context.Context,
		tenant *model.Tenant,
	) (*model.Tenant, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}

type ListTenantsQuery struct {
	BaseListOptions
}
