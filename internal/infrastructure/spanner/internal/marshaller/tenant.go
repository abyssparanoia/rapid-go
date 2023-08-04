package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/dbmodel"
)

func TenantToModel(e *dbmodel.Tenant) *model.Tenant {
	m := &model.Tenant{
		ID:        e.TenantID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
	return m
}

func TenantsToModel(slice dbmodel.TenantSlice) model.Tenants {
	dsts := make(model.Tenants, len(slice))
	for idx, e := range slice {
		dsts[idx] = TenantToModel(e)
	}
	return dsts
}

func TenantToDBModel(m *model.Tenant) *dbmodel.Tenant {
	e := &dbmodel.Tenant{}
	e.TenantID = m.ID
	e.Name = m.Name
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt

	return e
}
