package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/dbmodel"
	"github.com/volatiletech/null/v8"
)

func StaffToModel(e *dbmodel.Staff) *model.Staff {
	m := &model.Staff{
		ID:          e.StaffID,
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
		StaffID:     m.ID,
		TenantID:    m.TenantID,
		Role:        m.Role.String(),
		AuthUID:     m.AuthUID,
		DisplayName: m.DisplayName,
		ImagePath:   m.ImagePath,
		Email:       m.Email,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
