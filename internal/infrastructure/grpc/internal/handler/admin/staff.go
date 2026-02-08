package admin

import (
	"context"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) GetStaff(ctx context.Context, req *admin_apiv1.GetStaffRequest) (*admin_apiv1.GetStaffResponse, error) {
	got, err := h.staffInteractor.Get(
		ctx,
		input.NewAdminGetStaff(
			req.GetStaffId(),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}

	return &admin_apiv1.GetStaffResponse{
		Staff: marshaller.StaffToPB(got),
	}, nil
}

func (h *AdminHandler) ListStaffs(ctx context.Context, req *admin_apiv1.ListStaffsRequest) (*admin_apiv1.ListStaffsResponse, error) {
	var sortKey nullable.Type[model.StaffSortKey]
	if req.SortKey != nil {
		sortKey = nullable.TypeFrom(marshaller.StaffSortKeyToModel(*req.SortKey))
	}

	got, err := h.staffInteractor.List(
		ctx,
		input.NewAdminListStaffs(
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

	return &admin_apiv1.ListStaffsResponse{
		Staffs:     marshaller.StaffsToPB(got.Staffs),
		Pagination: marshaller.NewPagination(got.Pagination),
	}, nil
}

func (h *AdminHandler) CreateStaff(ctx context.Context, req *admin_apiv1.CreateStaffRequest) (*admin_apiv1.CreateStaffResponse, error) {
	claims, err := session_interceptor.RequireAdminSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	got, err := h.staffInteractor.Create(
		ctx,
		input.NewAdminCreateStaff(
			claims.AdminID.String,
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

	return &admin_apiv1.CreateStaffResponse{
		Staff:    marshaller.StaffToPB(got.Staff),
		Password: got.Password,
	}, nil
}

func (h *AdminHandler) UpdateStaff(ctx context.Context, req *admin_apiv1.UpdateStaffRequest) (*admin_apiv1.UpdateStaffResponse, error) {
	claims, err := session_interceptor.RequireAdminSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	var role nullable.Type[model.StaffRole]
	if req.Role != nil {
		role = nullable.TypeFrom(marshaller.StaffRoleToModel(*req.Role))
	}

	got, err := h.staffInteractor.Update(ctx, input.NewAdminUpdateStaff(
		claims.AdminID.String,
		req.GetStaffId(),
		null.StringFromPtr(req.DisplayName),
		role,
		null.StringFromPtr(req.ImageAssetId),
		request_interceptor.GetRequestTime(ctx),
	))
	if err != nil {
		return nil, err
	}

	return &admin_apiv1.UpdateStaffResponse{
		Staff: marshaller.StaffToPB(got),
	}, nil
}
