package staff

import (
	"context"

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
			claims.TenantID.String,
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
			claims.TenantID.String,
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
