package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/marshaller"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) AdminCreateStaff(ctx context.Context, req *admin_apiv1.AdminCreateStaffRequest) (*admin_apiv1.AdminCreateStaffResponse, error) {
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
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}

	return &admin_apiv1.AdminCreateStaffResponse{
		Staff: marshaller.StaffToPB(got),
	}, nil
}
