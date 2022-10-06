package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
	"github.com/volatiletech/null/v8"
)

type adminTenantInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
}

func NewAdminTenantInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
) AdminTenantInteractor {
	return &adminTenantInteractor{
		transactable,
		tenantRepository,
	}
}

func (i *adminTenantInteractor) Get(
	ctx context.Context,
	param *input.AdminGetTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	return i.tenantRepository.Get(ctx, param.TenantID, true)
}

func (i *adminTenantInteractor) List(
	ctx context.Context,
	param *input.AdminListTenants,
) (*output.ListTenants, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	tenants, err := i.tenantRepository.List(
		ctx,
		repository.ListTenantsQuery{
			Page:  null.Uint64From(param.Page),
			Limit: null.Uint64From(param.Limit),
		},
	)
	if err != nil {
		return nil, err
	}
	ttl, err := i.tenantRepository.Count(
		ctx,
		repository.CountTenantsQuery{},
	)
	if err != nil {
		return nil, err
	}
	return output.NewAdminListTenants(
		tenants,
		model.NewPagination(
			param.Page,
			param.Limit,
			ttl,
		),
	), nil
}

func (i *adminTenantInteractor) Create(
	ctx context.Context,
	param *input.AdminCreateTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	tenant := model.NewTenant(param.Name, param.RequestTime)
	return i.tenantRepository.Create(ctx, tenant)
}

func (i *adminTenantInteractor) Update(
	ctx context.Context,
	param *input.AdminUpdateTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	var tenant *model.Tenant
	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		var err error
		tenant, err = i.tenantRepository.Get(ctx, param.TenantID, true)
		if err != nil {
			return err
		}
		tenant.Update(
			null.StringFrom(param.Name),
			param.RequestTime,
		)
		tenant, err = i.tenantRepository.Update(ctx, tenant)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (i *adminTenantInteractor) Delete(
	ctx context.Context,
	param *input.AdminDeleteTenant,
) error {
	if err := param.Validate(); err != nil {
		return err
	}
	return i.tenantRepository.Delete(ctx, param.TenantID)
}
