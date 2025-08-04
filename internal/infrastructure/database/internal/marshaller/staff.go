package marshaller

import (
	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/dbmodel"
)

func StaffToModel(e *dbmodel.Staff) *model.Staff {
	m := &model.Staff{
		ID:          e.ID,
		TenantID:    e.TenantID,
		Role:        model.NewStaffRole(e.Role),
		AuthUID:     e.AuthUID,
		DisplayName: e.DisplayName,
		ImagePath:   e.ImagePath,
		Email:       e.Email,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,

		ImageURL:          null.String{},
		ReadonlyReference: nil,
	}

	if e.R != nil {
		var tenant *model.Tenant
		if e.R.Tenant != nil {
			tenant = TenantToModel(e.R.Tenant)
		}
		if tenant == nil {
			return m
		}
		m.ReadonlyReference = &struct {
			Tenant *model.Tenant
		}{
			Tenant: tenant,
		}
	}

	return m
}

func StaffsToModel(slice dbmodel.StaffSlice) model.Staffs {
	dsts := make(model.Staffs, len(slice))
	for idx, e := range slice {
		dsts[idx] = StaffToModel(e)
	}
	return dsts
}

func StaffToDBModel(m *model.Staff) *dbmodel.Staff {
	return &dbmodel.Staff{
		ID:          m.ID,
		TenantID:    m.TenantID,
		Role:        m.Role.String(),
		AuthUID:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImagePath:   m.ImagePath,
		Email:       m.Email,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		// R and L are likely relationship fields, initialize as nil if not needed
		R: nil,
		L: struct{}{},
	}
}
