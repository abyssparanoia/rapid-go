package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/volatiletech/null/v8"
)

type adminStaffInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	staffService     service.Staff
}

func NewAdminStaffInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	staffService service.Staff,
) AdminStaffInteractor {
	return &adminStaffInteractor{
		transactable,
		tenantRepository,
		staffService,
	}
}

func (i *adminStaffInteractor) Create(
	ctx context.Context,
	param *input.AdminCreateStaff,
) (*model.Staff, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	var staff *model.Staff
	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
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
			return err
		}

		staff, err = i.staffService.Create(
			ctx,
			service.StaffCreateParam{
				TenantID:    tenant.ID,
				Email:       param.Email,
				Password:    "random1234",
				StaffRole:   param.Role,
				DisplayName: param.DisplayName,
				ImagePath:   "staff_profile_images/default_image.jpeg",
				RequestTime: param.RequestTime,
			},
		)
		if err != nil {
			return err
		}

		staff.Tenant = tenant

		return nil
	}); err != nil {
		return nil, err
	}

	return staff, nil
}
