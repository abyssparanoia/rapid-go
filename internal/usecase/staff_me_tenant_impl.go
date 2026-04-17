package usecase

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type staffMeTenantInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	assetService     service.Asset
}

func NewStaffMeTenantInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	assetService service.Asset,
) StaffMeTenantInteractor {
	return &staffMeTenantInteractor{
		transactable:     transactable,
		tenantRepository: tenantRepository,
		assetService:     assetService,
	}
}

func (i *staffMeTenantInteractor) Get(
	ctx context.Context,
	param *input.StaffGetMeTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	tenant, err := i.tenantRepository.Get(ctx, repository.GetTenantQuery{
		ID: null.StringFrom(param.TenantID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// Apply asset URL processing
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}, param.RequestTime); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (i *staffMeTenantInteractor) Update(
	ctx context.Context,
	param *input.StaffUpdateMeTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		// Get with lock
		tenant, err := i.tenantRepository.Get(ctx, repository.GetTenantQuery{
			ID: null.StringFrom(param.TenantID),
			BaseGetOptions: repository.BaseGetOptions{
				OrFail:    true,
				ForUpdate: true,
			},
		})
		if err != nil {
			return err
		}

		// Apply updates via domain method
		tenant.Update(param.Name, param.RequestTime)

		// Persist
		return i.tenantRepository.Update(ctx, tenant)
	}); err != nil {
		return nil, err
	}

	// Return updated entity with relations
	tenant, err := i.tenantRepository.Get(ctx, repository.GetTenantQuery{
		ID: null.StringFrom(param.TenantID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// Apply asset URL processing
	if err := i.assetService.BatchSetTenantURLs(ctx, model.Tenants{tenant}, param.RequestTime); err != nil {
		return nil, err
	}

	return tenant, nil
}
