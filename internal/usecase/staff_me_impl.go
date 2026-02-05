package usecase

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

type staffMeInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	staffRepository  repository.Staff
	staffService     service.Staff
	assetService     service.Asset
}

func NewStaffMeInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	staffRepository repository.Staff,
	staffService service.Staff,
	assetService service.Asset,
) StaffMeInteractor {
	return &staffMeInteractor{
		transactable:     transactable,
		tenantRepository: tenantRepository,
		staffRepository:  staffRepository,
		staffService:     staffService,
		assetService:     assetService,
	}
}

func (i *staffMeInteractor) SignUp(
	ctx context.Context,
	param *input.StaffSignUp,
) (*model.Staff, error) {
	// 1. Validate input
	if err := param.Validate(); err != nil {
		return nil, err
	}

	// 2. Validate asset
	imagePath, err := i.assetService.GetWithValidate(ctx, model.AssetTypeUserImage, param.ImageAssetID)
	if err != nil {
		return nil, err
	}

	var staff *model.Staff

	// 3. Create tenant and staff in transaction
	txErr := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		// Create new tenant
		tenant := model.NewTenant(
			param.TenantName,
			param.RequestTime,
		)
		if createErr := i.tenantRepository.Create(ctx, tenant); createErr != nil {
			return createErr
		}

		// Create staff with admin role using domain service
		var createStaffErr error
		staff, createStaffErr = i.staffService.Create(ctx, service.StaffCreateParam{
			TenantID:    tenant.ID,
			Email:       param.Email,
			Password:    "", // User already has Cognito account
			StaffRole:   model.StaffRoleAdmin,
			DisplayName: param.DisplayName,
			ImagePath:   imagePath,
			RequestTime: param.RequestTime,
		})
		if createStaffErr != nil {
			return createStaffErr
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	// 4. Return staff with relations loaded
	staff, err = i.staffRepository.Get(ctx, repository.GetStaffQuery{
		ID: null.StringFrom(staff.ID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// 5. Apply asset URL processing
	if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}); err != nil {
		return nil, err
	}

	return staff, nil
}

func (i *staffMeInteractor) Get(
	ctx context.Context,
	param *input.StaffGetMe,
) (*model.Staff, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
		ID: null.StringFrom(param.StaffID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// Apply asset URL processing
	if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}); err != nil {
		return nil, err
	}

	return staff, nil
}

func (i *staffMeInteractor) Update(
	ctx context.Context,
	param *input.StaffUpdateMe,
) (*model.Staff, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	// Validate asset if provided
	var imagePath string
	if param.ImageAssetID.Valid {
		var err error
		imagePath, err = i.assetService.GetWithValidate(ctx, model.AssetTypeUserImage, param.ImageAssetID.String)
		if err != nil {
			return nil, err
		}
	}

	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		// Get with lock
		staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
			ID: null.StringFrom(param.StaffID),
			BaseGetOptions: repository.BaseGetOptions{
				OrFail:    true,
				ForUpdate: true,
			},
		})
		if err != nil {
			return err
		}

		// Apply updates via domain method (role is not updated in UpdateMe)
		staff.Update(param.DisplayName, nullable.Type[model.StaffRole]{}, null.StringFrom(imagePath), param.RequestTime)

		// Persist
		return i.staffRepository.Update(ctx, staff)
	}); err != nil {
		return nil, err
	}

	// Return updated entity with relations
	staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
		ID: null.StringFrom(param.StaffID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// Apply asset URL processing
	if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}); err != nil {
		return nil, err
	}

	return staff, nil
}
