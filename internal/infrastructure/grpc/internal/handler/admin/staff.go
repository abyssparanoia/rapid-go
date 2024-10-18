package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/admin/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) CreateStaff(ctx context.Context, req *admin_apiv1.CreateStaffRequest) (*admin_apiv1.CreateStaffResponse, error) {
	_, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	got, err := h.staffInteractor.Create(
		ctx,
		input.NewAdminCreateStaff(
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
		Staff: marshaller.StaffToPB(got),
	}, nil
}
