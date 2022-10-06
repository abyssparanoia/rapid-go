package marshaller

import (
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/database/internal/dbmodel"
)

func OutputTenantToModel(e *dbmodel.Tenant) *model.Tenant {
	m := &model.Tenant{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
	return m
}

func OutputTenantsToModel(slice dbmodel.TenantSlice) []*model.Tenant {
	dsts := make([]*model.Tenant, len(slice))
	for idx, e := range slice {
		dsts[idx] = OutputTenantToModel(e)
	}
	return dsts
}

func NewTenantFromModel(m *model.Tenant) *dbmodel.Tenant {
	e := &dbmodel.Tenant{}
	e.ID = m.ID
	e.Name = m.Name
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt

	return e
}
