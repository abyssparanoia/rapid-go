package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type AdminTenantInteractor interface {
	Get(
		ctx context.Context,
		param *input.AdminGetTenant,
	) (*model.Tenant, error)
	List(
		ctx context.Context,
		param *input.AdminListTenants,
	) (*output.ListTenants, error)
	Create(
		ctx context.Context,
		param *input.AdminCreateTenant,
	) (*model.Tenant, error)
	Update(
		ctx context.Context,
		param *input.AdminUpdateTenant,
	) (*model.Tenant, error)
	Delete(
		ctx context.Context,
		param *input.AdminDeleteTenant,
	) error
}
