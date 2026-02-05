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

type staffTenantInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	assetService     service.Asset
}

func NewStaffTenantInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	assetService service.Asset,
) StaffTenantInteractor {
	return &staffTenantInteractor{
		transactable,
		tenantRepository,
		assetService,
	}
}

func (i *staffTenantInteractor) Get(
	ctx context.Context,
	param *input.StaffGetTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	tenant, err := i.tenantRepository.Get(
		ctx,
		repository.GetTenantQuery{
			ID: null.StringFrom(param.TenantID),
			BaseGetOptions: repository.BaseGetOptions{
				OrFail: true,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (i *staffTenantInteractor) List(
	ctx context.Context,
	param *input.StaffListTenants,
) (*output.ListTenants, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	query := repository.ListTenantsQuery{
		BaseListOptions: repository.BaseListOptions{
			Page:  null.Uint64From(param.Page),
			Limit: null.Uint64From(param.Limit),
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
	if err = i.assetService.BatchSetTenantURLs(ctx, tenants); err != nil {
		return nil, err
	}
	ttl, err := i.tenantRepository.Count(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	return output.NewStaffListTenants(
		tenants,
		model.NewPagination(
			param.Page,
			param.Limit,
			ttl,
		),
	), nil
}

func (i *staffTenantInteractor) Create(
	ctx context.Context,
	param *input.StaffCreateTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	tenant := model.NewTenant(param.Name, param.RequestTime)
	if err := i.tenantRepository.Create(ctx, tenant); err != nil {
		return nil, err
	}
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (i *staffTenantInteractor) Update(
	ctx context.Context,
	param *input.StaffUpdateTenant,
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
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (i *staffTenantInteractor) Delete(
	ctx context.Context,
	param *input.StaffDeleteTenant,
) error {
	if err := param.Validate(); err != nil {
		return err
	}
	return i.tenantRepository.Delete(ctx, param.TenantID)
}
