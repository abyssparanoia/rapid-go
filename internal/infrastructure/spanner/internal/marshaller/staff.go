package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/dbmodel"
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
	e := &dbmodel.Staff{}
	e.StaffID = m.ID
	e.TenantID = m.TenantID
	e.Role = m.Role.String()
	e.AuthUID = m.AuthUID
	e.DisplayName = m.DisplayName
	e.ImagePath = m.ImagePath
	e.Email = m.Email
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt

	return e
}
