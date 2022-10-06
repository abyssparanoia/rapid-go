package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/marshaller"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *AdminHandler) GetTenant(ctx context.Context, req *admin_apiv1.GetTenantRequest) (*admin_apiv1.GetTenantResponse, error) {
	got, err := h.tenantInteractor.Get(
		ctx,
		input.NewAdminGetTenant(
			req.GetTenantId(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.GetTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *AdminHandler) ListTenants(ctx context.Context, req *admin_apiv1.ListTenantsRequest) (*admin_apiv1.ListTenantsResponse, error) {
	got, err := h.tenantInteractor.List(
		ctx,
		input.NewAdminListTenants(
			req.GetPage(),
			req.GetLimit(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.ListTenantsResponse{
		Tenants:    marshaller.TenantsToPB(got.Tenants),
		Pagination: marshaller.NewPagination(got.Pagination),
	}, nil
}

func (h *AdminHandler) CreateTenant(ctx context.Context, req *admin_apiv1.CreateTenantRequest) (*admin_apiv1.CreateTenantResponse, error) {
	got, err := h.tenantInteractor.Create(
		ctx,
		input.NewAdminCreateTenant(
			req.GetName(),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.CreateTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *AdminHandler) UpdateTenant(ctx context.Context, req *admin_apiv1.UpdateTenantRequest) (*admin_apiv1.UpdateTenantResponse, error) {
	got, err := h.tenantInteractor.Update(
		ctx,
		input.NewAdminUpdateTenant(
			req.GetTenantId(),
			req.GetName(),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.UpdateTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *AdminHandler) DeleteTenant(ctx context.Context, req *admin_apiv1.DeleteTenantRequest) (*admin_apiv1.DeleteTenantResponse, error) {
	err := h.tenantInteractor.Delete(
		ctx,
		input.NewAdminDeleteTenant(
			req.GetTenantId(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.DeleteTenantResponse{}, nil
}
