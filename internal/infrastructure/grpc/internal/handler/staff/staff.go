package staff

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/staff/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *StaffHandler) GetStaff(ctx context.Context, req *staff_apiv1.GetStaffRequest) (*staff_apiv1.GetStaffResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	got, err := h.staffInteractor.Get(
		ctx,
		input.NewStaffGetStaff(
			claims.StaffID.String,
			req.GetStaffId(),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.GetStaffResponse{
		Staff: marshaller.StaffToPB(got),
	}, nil
}

func (h *StaffHandler) ListStaffs(ctx context.Context, req *staff_apiv1.ListStaffsRequest) (*staff_apiv1.ListStaffsResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	var sortKey nullable.Type[model.StaffSortKey]
	if req.SortKey != nil {
		sortKey = nullable.TypeFrom(marshaller.StaffSortKeyToModel(*req.SortKey))
	}

	got, err := h.staffInteractor.List(
		ctx,
		input.NewStaffListStaffs(
			claims.StaffID.String,
			req.GetTenantId(),
			req.GetPage(),
			req.GetLimit(),
			sortKey,
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.ListStaffsResponse{
		Staffs:     marshaller.StaffsToPB(got.Staffs),
		Pagination: marshaller.NewPagination(got.Pagination),
	}, nil
}

func (h *StaffHandler) CreateStaff(ctx context.Context, req *staff_apiv1.CreateStaffRequest) (*staff_apiv1.CreateStaffResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	got, err := h.staffInteractor.Create(
		ctx,
		input.NewStaffCreateStaff(
			claims.StaffID.String,
			req.GetTenantId(),
			req.GetEmail(),
			req.GetDisplayName(),
			marshaller.StaffRoleToModel(req.GetRole()),
			req.GetImageAssetId(),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.CreateStaffResponse{
		Staff: marshaller.StaffToPB(got),
	}, nil
}

func (h *StaffHandler) UpdateStaff(ctx context.Context, req *staff_apiv1.UpdateStaffRequest) (*staff_apiv1.UpdateStaffResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	// Build input parameter with required fields
	param := input.NewStaffUpdateStaff(
		claims.StaffID.String,
		req.GetStaffId(),
		null.String{},
		nullable.Type[model.StaffRole]{},
		null.String{},
		request_interceptor.GetRequestTime(ctx),
	)

	// Set optional fields if provided
	if req.DisplayName != nil {
		param.DisplayName = null.StringFrom(*req.DisplayName)
	}
	if req.Role != nil {
		param.Role = nullable.TypeFrom(marshaller.StaffRoleToModel(*req.Role))
	}
	if req.ImageAssetId != nil {
		param.ImageAssetID = null.StringFrom(*req.ImageAssetId)
	}

	got, err := h.staffInteractor.Update(ctx, param)
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.UpdateStaffResponse{
		Staff: marshaller.StaffToPB(got),
	}, nil
}
