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

func (h *StaffHandler) GetTenant(ctx context.Context, req *staff_apiv1.GetTenantRequest) (*staff_apiv1.GetTenantResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	got, err := h.tenantInteractor.Get(
		ctx,
		input.NewStaffGetTenant(
			claims.StaffID.String,
			req.GetTenantId(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &staff_apiv1.GetTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *StaffHandler) ListTenants(ctx context.Context, req *staff_apiv1.ListTenantsRequest) (*staff_apiv1.ListTenantsResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	var sortKey nullable.Type[model.TenantSortKey]
	if req.SortKey != nil {
		sortKey = nullable.TypeFrom(marshaller.TenantSortKeyToModel(*req.SortKey))
	}

	got, err := h.tenantInteractor.List(
		ctx,
		input.NewStaffListTenants(
			claims.StaffID.String,
			req.GetPage(),
			req.GetLimit(),
			sortKey,
		),
	)
	if err != nil {
		return nil, err
	}
	return &staff_apiv1.ListTenantsResponse{
		Tenants:    marshaller.TenantsToPB(got.Tenants),
		Pagination: marshaller.NewPagination(got.Pagination),
	}, nil
}

func (h *StaffHandler) CreateTenant(ctx context.Context, req *staff_apiv1.CreateTenantRequest) (*staff_apiv1.CreateTenantResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	got, err := h.tenantInteractor.Create(
		ctx,
		input.NewStaffCreateTenant(
			claims.StaffID.String,
			req.GetName(),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}
	return &staff_apiv1.CreateTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *StaffHandler) UpdateTenant(ctx context.Context, req *staff_apiv1.UpdateTenantRequest) (*staff_apiv1.UpdateTenantResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	got, err := h.tenantInteractor.Update(
		ctx,
		input.NewStaffUpdateTenant(
			claims.StaffID.String,
			req.GetTenantId(),
			null.StringFromPtr(req.Name),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}
	return &staff_apiv1.UpdateTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *StaffHandler) DeleteTenant(ctx context.Context, req *staff_apiv1.DeleteTenantRequest) (*staff_apiv1.DeleteTenantResponse, error) {
	claims, err := session_interceptor.RequireStaffSessionContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.tenantInteractor.Delete(
		ctx,
		input.NewStaffDeleteTenant(
			claims.StaffID.String,
			req.GetTenantId(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &staff_apiv1.DeleteTenantResponse{}, nil
}
