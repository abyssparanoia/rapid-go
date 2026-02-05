package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type StaffMeTenantInteractor interface {
	Get(
		ctx context.Context,
		param *input.StaffGetMeTenant,
	) (*model.Tenant, error)
	Update(
		ctx context.Context,
		param *input.StaffUpdateMeTenant,
	) (*model.Tenant, error)
}
