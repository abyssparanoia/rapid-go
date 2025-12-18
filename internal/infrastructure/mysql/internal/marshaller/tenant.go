package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/internal/dbmodel"
)

func TenantToModel(e *dbmodel.Tenant) *model.Tenant {
	m := &model.Tenant{
		ID:        e.ID,
		Name:      e.Name,
		Tags:      make(model.TenantTags, 0),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	if e.R != nil && e.R.TenantTags != nil {
		m.Tags = TenantTagsToModel(e.R.TenantTags)
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
	return &dbmodel.Tenant{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		// Initialize R and L if they exist in dbmodel.Tenant
		R: nil,
		L: struct{}{},
	}
}

func TenantsToDBModel(m model.Tenants) dbmodel.TenantSlice {
	dsts := make(dbmodel.TenantSlice, len(m))
	for idx, e := range m {
		dsts[idx] = TenantToDBModel(e)
	}
	return dsts
}

func TenantTagToModel(e *dbmodel.TenantTag) *model.TenantTag {
	return &model.TenantTag{
		ID:        e.ID,
		Type:      model.NewTenantTagType(e.Type),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func TenantTagsToModel(slice dbmodel.TenantTagSlice) model.TenantTags {
	dsts := make(model.TenantTags, len(slice))
	for idx, e := range slice {
		dsts[idx] = TenantTagToModel(e)
	}
	return dsts
}

func TenantTagToDBModel(m *model.TenantTag, tenantID string) *dbmodel.TenantTag {
	return &dbmodel.TenantTag{
		ID:        m.ID,
		TenantID:  tenantID,
		Type:      m.Type.String(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		R:         nil,
		L:         struct{}{},
	}
}

func TenantTagsToDBModel(m model.TenantTags, tenantID string) dbmodel.TenantTagSlice {
	dsts := make(dbmodel.TenantTagSlice, len(m))
	for idx, e := range m {
		dsts[idx] = TenantTagToDBModel(e, tenantID)
	}
	return dsts
}
