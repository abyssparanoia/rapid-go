package usecase

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type adminStaffInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	staffRepository  repository.Staff
	staffService     service.Staff
	assetService     service.Asset
}

func NewAdminStaffInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	staffRepository repository.Staff,
	staffService service.Staff,
	assetService service.Asset,
) AdminStaffInteractor {
	return &adminStaffInteractor{
		transactable,
		tenantRepository,
		staffRepository,
		staffService,
		assetService,
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

		imagePath, err := i.assetService.GetWithValidate(ctx, model.AssetTypeUserImage, param.ImageAssetID)
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
				ImagePath:   imagePath,
				RequestTime: param.RequestTime,
			},
		)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// 表示用に再取得する
	staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
		ID: null.StringFrom(staff.ID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return staff, nil
}
