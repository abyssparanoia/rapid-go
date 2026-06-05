package usecase

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type adminTenantInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	assetService     service.Asset
}

func NewAdminTenantInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	assetService service.Asset,
) AdminTenantInteractor {
	return &adminTenantInteractor{
		transactable,
		tenantRepository,
		assetService,
	}
}

func (i *adminTenantInteractor) Get(
	ctx context.Context,
	param *input.AdminGetTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	tenant, err := i.tenantRepository.Get(
		ctx,
		repository.GetTenantQuery{
			BaseGetOptions: repository.BaseGetOptions{
				OrFail:  true,
				Preload: true,
			},
			ID: null.StringFrom(param.TenantID),
		},
	)
	if err != nil {
		return nil, err
	}
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}, param.RequestTime); err != nil {
		return nil, err
	}
	return tenant, nil
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
			Page:    null.Uint64From(param.Page),
			Limit:   null.Uint64From(param.Limit),
			Preload: true,
		},
		SortKey: nullable.TypeFrom(param.SortKey),
	}
	tenants, err := i.tenantRepository.List(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	if err = i.assetService.BatchSetTenantURLs(ctx, tenants, param.RequestTime); err != nil {
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
	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		return i.tenantRepository.Create(ctx, tenant)
	}); err != nil {
		return nil, err
	}
	tenant, err := i.tenantRepository.Get(ctx, repository.GetTenantQuery{
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
		ID: null.StringFrom(tenant.ID),
	})
	if err != nil {
		return nil, err
	}
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}, param.RequestTime); err != nil {
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
	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		tenant, err := i.tenantRepository.Get(
			ctx,
			repository.GetTenantQuery{
				BaseGetOptions: repository.BaseGetOptions{
					OrFail:    true,
					ForUpdate: true,
				},
				ID: null.StringFrom(param.TenantID),
			},
		)
		if err != nil {
			return err
		}
		tenant.Update(
			param.Name,
			param.RequestTime,
		)
		return i.tenantRepository.Update(ctx, tenant)
	}); err != nil {
		return nil, err
	}
	tenant, err := i.tenantRepository.Get(ctx, repository.GetTenantQuery{
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
		ID: null.StringFrom(param.TenantID),
	})
	if err != nil {
		return nil, err
	}
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}, param.RequestTime); err != nil {
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
	return i.transactable.RWTx(ctx, func(ctx context.Context) error {
		_, err := i.tenantRepository.Get(ctx, repository.GetTenantQuery{
			BaseGetOptions: repository.BaseGetOptions{
				OrFail:    true,
				ForUpdate: true,
			},
			ID: null.StringFrom(param.TenantID),
		})
		if err != nil {
			return err
		}
		return i.tenantRepository.Delete(ctx, param.TenantID)
	})
}
