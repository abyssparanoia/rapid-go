package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type staffInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	staffService     service.Staff
}

func NewStaffInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	staffService service.Staff,
) StaffInteractor {
	return &staffInteractor{
		transactable,
		tenantRepository,
		staffService,
	}
}

func (i *staffInteractor) CreateRoot(
	ctx context.Context,
	param *input.CreateRootStaff,
) error {
	return i.transactable.RWTx(ctx, func(ctx context.Context) error {
		if err := param.Validate(); err != nil {
			return err
		}

		tenant := model.NewTenant("Platformer", param.RequestTime)
		if err := i.tenantRepository.Create(ctx, tenant); err != nil {
			return err
		}

		if _, err := i.staffService.Create(
			ctx,
			service.StaffCreateParam{
				TenantID:    tenant.ID,
				Email:       param.Email,
				Password:    "random1234",
				StaffRole:   model.StaffRoleAdmin,
				DisplayName: "Root Staff",
				ImagePath:   "staff_profile_images/default_image.jpeg",
				RequestTime: param.RequestTime,
			},
		); err != nil {
			return err
		}

		return nil
	})
}
