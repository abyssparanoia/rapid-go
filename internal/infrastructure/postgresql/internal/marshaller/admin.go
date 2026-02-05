package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/postgresql/internal/dbmodel"
)

func AdminToModel(e *dbmodel.Admin) *model.Admin {
	return &model.Admin{
		ID:          e.ID,
		Role:        model.NewAdminRole(e.Role),
		AuthUID:     e.AuthUID,
		Email:       e.Email,
		DisplayName: e.DisplayName,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func AdminsToModel(slice dbmodel.AdminSlice) model.Admins {
	dsts := make(model.Admins, len(slice))
	for idx, e := range slice {
		dsts[idx] = AdminToModel(e)
	}
	return dsts
}

func AdminToDBModel(m *model.Admin) *dbmodel.Admin {
	return &dbmodel.Admin{
		ID:          m.ID,
		Role:        m.Role.String(),
		AuthUID:     m.AuthUID,
		Email:       m.Email,
		DisplayName: m.DisplayName,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		R:           nil,
		L:           struct{}{},
	}
}

func AdminsToDBModel(m model.Admins) dbmodel.AdminSlice {
	dsts := make(dbmodel.AdminSlice, len(m))
	for idx, e := range m {
		dsts[idx] = AdminToDBModel(e)
	}
	return dsts
}
