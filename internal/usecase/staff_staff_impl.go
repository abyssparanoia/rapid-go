package usecase

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type staffStaffInteractor struct {
	transactable     repository.Transactable
	tenantRepository repository.Tenant
	staffRepository  repository.Staff
	staffService     service.Staff
	assetService     service.Asset
}

func NewStaffStaffInteractor(
	transactable repository.Transactable,
	tenantRepository repository.Tenant,
	staffRepository repository.Staff,
	staffService service.Staff,
	assetService service.Asset,
) StaffStaffInteractor {
	return &staffStaffInteractor{
		transactable,
		tenantRepository,
		staffRepository,
		staffService,
		assetService,
	}
}

func (i *staffStaffInteractor) Get(
	ctx context.Context,
	param *input.StaffGetStaff,
) (*model.Staff, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
		ID: null.StringFrom(param.TargetStaffID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}); err != nil {
		return nil, err
	}

	return staff, nil
}

func (i *staffStaffInteractor) List(
	ctx context.Context,
	param *input.StaffListStaffs,
) (*output.ListStaffs, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	query := repository.ListStaffQuery{
		TenantID: null.StringFrom(param.TenantID),
		BaseListOptions: repository.BaseListOptions{
			Page:    null.Uint64From(param.Page),
			Limit:   null.Uint64From(param.Limit),
			Preload: true,
		},
		SortKey: nullable.TypeFrom(param.SortKey),
	}

	staffs, err := i.staffRepository.List(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = i.assetService.BatchSetStaffURLs(ctx, staffs); err != nil {
		return nil, err
	}

	totalCount, err := i.staffRepository.Count(ctx, query)
	if err != nil {
		return nil, err
	}

	return output.NewStaffListStaffs(
		staffs,
		model.NewPagination(
			param.Page,
			param.Limit,
			totalCount,
		),
	), nil
}

func (i *staffStaffInteractor) Create(
	ctx context.Context,
	param *input.StaffCreateStaff,
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

	if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}); err != nil {
		return nil, err
	}

	return staff, nil
}

func (i *staffStaffInteractor) Update(
	ctx context.Context,
	param *input.StaffUpdateStaff,
) (*model.Staff, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}

	var staff *model.Staff
	if err := i.transactable.RWTx(ctx, func(ctx context.Context) error {
		var err error
		staff, err = i.staffRepository.Get(ctx, repository.GetStaffQuery{
			ID: null.StringFrom(param.TargetStaffID),
			BaseGetOptions: repository.BaseGetOptions{
				OrFail:    true,
				ForUpdate: true,
			},
		})
		if err != nil {
			return err
		}

		// Get image path if image_asset_id is provided
		var imagePath null.String
		if param.ImageAssetID.Valid {
			path, err := i.assetService.GetWithValidate(ctx, model.AssetTypeUserImage, param.ImageAssetID.String)
			if err != nil {
				return err
			}
			imagePath = null.StringFrom(path)
		}

		// Update via domain method
		staff.Update(
			param.DisplayName,
			param.Role,
			imagePath,
			param.RequestTime,
		)

		if err := i.staffRepository.Update(ctx, staff); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Retrieve with preload for response
	staff, err := i.staffRepository.Get(ctx, repository.GetStaffQuery{
		ID: null.StringFrom(param.TargetStaffID),
		BaseGetOptions: repository.BaseGetOptions{
			OrFail:  true,
			Preload: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := i.assetService.BatchSetStaffURLs(ctx, model.Staffs{staff}); err != nil {
		return nil, err
	}

	return staff, nil
}
