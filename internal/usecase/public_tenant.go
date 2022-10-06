package usecase

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
)

type PublicTenantInteractor interface {
	Get(
		ctx context.Context,
		param *input.PublicGetTenant,
	) (*model.Tenant, error)
}
