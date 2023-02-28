package admin

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/request_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/marshaller"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/volatiletech/null/v8"
)

func (h *AdminHandler) AdminGetTenant(ctx context.Context, req *admin_apiv1.AdminGetTenantRequest) (*admin_apiv1.AdminGetTenantResponse, error) {
	got, err := h.tenantInteractor.Get(
		ctx,
		input.NewAdminGetTenant(
			req.GetTenantId(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.AdminGetTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *AdminHandler) AdminListTenants(ctx context.Context, req *admin_apiv1.AdminListTenantsRequest) (*admin_apiv1.AdminListTenantsResponse, error) {
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
	return &admin_apiv1.AdminListTenantsResponse{
		Tenants:    marshaller.TenantsToPB(got.Tenants),
		Pagination: marshaller.NewPagination(got.Pagination),
	}, nil
}

func (h *AdminHandler) AdminCreateTenant(ctx context.Context, req *admin_apiv1.AdminCreateTenantRequest) (*admin_apiv1.AdminCreateTenantResponse, error) {
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
	return &admin_apiv1.AdminCreateTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *AdminHandler) AdminUpdateTenant(ctx context.Context, req *admin_apiv1.AdminUpdateTenantRequest) (*admin_apiv1.AdminUpdateTenantResponse, error) {
	got, err := h.tenantInteractor.Update(
		ctx,
		input.NewAdminUpdateTenant(
			req.GetTenantId(),
			null.StringFromPtr(req.Name),
			request_interceptor.GetRequestTime(ctx),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.AdminUpdateTenantResponse{
		Tenant: marshaller.TenantToPB(got),
	}, nil
}

func (h *AdminHandler) AdminDeleteTenant(ctx context.Context, req *admin_apiv1.AdminDeleteTenantRequest) (*admin_apiv1.AdminDeleteTenantResponse, error) {
	err := h.tenantInteractor.Delete(
		ctx,
		input.NewAdminDeleteTenant(
			req.GetTenantId(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &admin_apiv1.AdminDeleteTenantResponse{}, nil
}
