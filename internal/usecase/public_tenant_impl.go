package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type publicTenantInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
}

func NewPublicTenantInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
) PublicTenantInteractor {
	return &publicTenantInteractor{
		transactable,
		tenantRepository,
	}
}

func (i *publicTenantInteractor) Get(
	ctx context.Context,
	param *input.PublicGetTenant,
) (*model.Tenant, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	return i.tenantRepository.Get(ctx, param.TenantID, true)
}
