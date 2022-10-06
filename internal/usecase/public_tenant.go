package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type PublicTenantInteractor interface {
	Get(
		ctx context.Context,
		param *input.PublicGetTenant,
	) (*model.Tenant, error)
}
