package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/volatiletech/null/v8"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Tenant interface {
	Get(
		ctx context.Context,
		query GetTenantQuery,
	) (*model.Tenant, error)
	List(
		ctx context.Context,
		query ListTenantsQuery,
	) (model.Tenants, error)
	Count(
		ctx context.Context,
		query ListTenantsQuery,
	) (uint64, error)
	Create(
		ctx context.Context,
		tenant *model.Tenant,
	) error
	Update(
		ctx context.Context,
		tenant *model.Tenant,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}

type GetTenantQuery struct {
	ID null.String
	BaseGetOptions
}

type ListTenantsQuery struct {
	BaseListOptions
}
