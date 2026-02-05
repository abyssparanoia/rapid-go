package staff

import (
	"context"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/handler/staff/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *StaffHandler) GetMeTenant(
	ctx context.Context,
	req *staff_apiv1.GetMeTenantRequest,
) (*staff_apiv1.GetMeTenantResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	requestTime := request_interceptor.GetRequestTime(ctx)

	tenant, err := h.meTenantInteractor.Get(
		ctx,
		input.NewStaffGetMeTenant(
			claims.TenantID.String,
			claims.StaffID.String,
			requestTime,
		),
	)
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.GetMeTenantResponse{
		Tenant: marshaller.TenantToPB(tenant),
	}, nil
}

func (h *StaffHandler) UpdateMeTenant(
	ctx context.Context,
	req *staff_apiv1.UpdateMeTenantRequest,
) (*staff_apiv1.UpdateMeTenantResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	requestTime := request_interceptor.GetRequestTime(ctx)

	param := input.NewStaffUpdateMeTenant(
		claims.TenantID.String,
		claims.StaffID.String,
		null.String{},
		requestTime,
	)

	if req.Name != nil {
		param.Name = null.StringFrom(*req.Name)
	}

	tenant, err := h.meTenantInteractor.Update(ctx, param)
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.UpdateMeTenantResponse{
		Tenant: marshaller.TenantToPB(tenant),
	}, nil
}
