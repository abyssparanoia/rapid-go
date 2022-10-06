package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/marshaller"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) CreateUser(ctx context.Context, req *admin_apiv1.CreateUserRequest) (*admin_apiv1.CreateUserResponse, error) {
	_, err := session_interceptor.RequireSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	got, err := h.userInteractor.Create(
		ctx,
		input.NewAdminCreateUser(
			req.GetTenantId(),
			req.GetEmail(),
			req.GetDisplayName(),
			marshaller.UserRoleToModel(req.GetRole()),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}

	return &admin_apiv1.CreateUserResponse{
		User: marshaller.UserToPB(got),
	}, nil
}
