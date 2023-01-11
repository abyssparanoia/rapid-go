package output

import "github.com/abyssparanoia/rapid-go/internal/domain/model"

type ListTenants struct {
	Tenants    model.Tenants
	Pagination *model.Pagination
}

func NewAdminListTenants(
	tenants model.Tenants,
	pagination *model.Pagination,
) *ListTenants {
	return &ListTenants{
		Tenants:    tenants,
		Pagination: pagination,
	}
}
