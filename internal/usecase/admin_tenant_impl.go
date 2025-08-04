package usecase

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
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
	return i.tenantRepository.Get(
		ctx,
		repository.GetTenantQuery{
			ID: null.StringFrom(param.TenantID),
			BaseGetOptions: repository.BaseGetOptions{
				OrFail: true,
			},
		},
	)
}

func (i *adminTenantInteractor) List(
	ctx context.Context,
	param *input.AdminListTenants,
) (*output.ListTenants, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	query := repository.ListTenantsQuery{
		BaseListOptions: repository.BaseListOptions{
			Page:  null.Uint64From(param.Page),
			Limit: null.Uint64From(param.Limit),
		},
	}
	tenants, err := i.tenantRepository.List(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	ttl, err := i.tenantRepository.Count(
		ctx,
		query,
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
	if err := i.tenantRepository.Create(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
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
		tenant, err = i.tenantRepository.Get(
			ctx,
			repository.GetTenantQuery{
				ID: null.StringFrom(param.TenantID),
				BaseGetOptions: repository.BaseGetOptions{
					OrFail:    true,
					ForUpdate: true,
				},
			},
		)
		if err != nil {
			return err
		}
		tenant.Update(
			param.Name,
			param.RequestTime,
		)
		if err := i.tenantRepository.Update(ctx, tenant); err != nil {
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
