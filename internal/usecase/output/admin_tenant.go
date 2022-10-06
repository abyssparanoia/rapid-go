package output

import "github.com/abyssparanoia/rapid-go/internal/domain/model"

type ListTenants struct {
	Tenants    []*model.Tenant
	Pagination *model.Pagination
}

func NewAdminListTenants(
	tenants []*model.Tenant,
	pagination *model.Pagination,
) *ListTenants {
	return &ListTenants{
		Tenants:    tenants,
		Pagination: pagination,
	}
}
