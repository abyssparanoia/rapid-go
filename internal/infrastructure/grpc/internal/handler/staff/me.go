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

func (h *StaffHandler) SignUp(
	ctx context.Context,
	req *staff_apiv1.SignUpRequest,
) (*staff_apiv1.SignUpResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	requestTime := request_interceptor.GetRequestTime(ctx)

	staff, err := h.meInteractor.SignUp(ctx, input.NewStaffSignUp(
		claims.AuthUID,
		claims.Email,
		req.GetTenantName(),
		req.GetDisplayName(),
		req.GetImageAssetId(),
		requestTime,
	))
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.SignUpResponse{
		Staff: marshaller.StaffToPB(staff),
	}, nil
}

func (h *StaffHandler) GetMe(
	ctx context.Context,
	req *staff_apiv1.GetMeRequest,
) (*staff_apiv1.GetMeResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	requestTime := request_interceptor.GetRequestTime(ctx)

	staff, err := h.meInteractor.Get(ctx, input.NewStaffGetMe(
		claims.TenantID.String,
		claims.StaffID.String,
		requestTime,
	))
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.GetMeResponse{
		Staff: marshaller.StaffToPB(staff),
	}, nil
}

func (h *StaffHandler) UpdateMe(
	ctx context.Context,
	req *staff_apiv1.UpdateMeRequest,
) (*staff_apiv1.UpdateMeResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	requestTime := request_interceptor.GetRequestTime(ctx)

	param := input.NewStaffUpdateMe(
		claims.TenantID.String,
		claims.StaffID.String,
		null.String{},
		null.String{},
		requestTime,
	)

	if req.DisplayName != nil {
		param.DisplayName = null.StringFrom(*req.DisplayName)
	}
	if req.ImageAssetId != nil {
		param.ImageAssetID = null.StringFrom(*req.ImageAssetId)
	}

	staff, err := h.meInteractor.Update(ctx, param)
	if err != nil {
		return nil, err
	}

	return &staff_apiv1.UpdateMeResponse{
		Staff: marshaller.StaffToPB(staff),
	}, nil
}
