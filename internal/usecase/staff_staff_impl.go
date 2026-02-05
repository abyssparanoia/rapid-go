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
