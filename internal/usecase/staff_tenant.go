package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type StaffTenantInteractor interface {
	Get(
		ctx context.Context,
		param *input.StaffGetTenant,
	) (*model.Tenant, error)
	List(
		ctx context.Context,
		param *input.StaffListTenants,
	) (*output.ListTenants, error)
	Create(
		ctx context.Context,
		param *input.StaffCreateTenant,
	) (*model.Tenant, error)
	Update(
		ctx context.Context,
		param *input.StaffUpdateTenant,
	) (*model.Tenant, error)
	Delete(
		ctx context.Context,
		param *input.StaffDeleteTenant,
	) error
}
