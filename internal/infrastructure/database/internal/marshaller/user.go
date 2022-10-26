package marshaller

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/internal/dbmodel"
)

func UserToModel(e *dbmodel.User) *model.User {
	m := &model.User{
		ID:          e.ID,
		Role:        model.NewUserRole(e.Role),
		AuthUID:     e.AuthUID,
		DisplayName: e.DisplayName,
		ImagePath:   e.ImagePath,
		Email:       e.Email,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	if e.R != nil && e.R.Tenant != nil {
		m.Tenant = TenantToModel(e.R.Tenant)
	} else {
		m.Tenant = &model.Tenant{
			ID: e.TenantID,
		}
	}

	return m
}

func UserToDBModel(m *model.User) *dbmodel.User {
	e := &dbmodel.User{}
	e.ID = m.ID
	e.TenantID = m.Tenant.ID
	e.Role = m.Role.String()
	e.AuthUID = m.AuthUID
	e.DisplayName = m.DisplayName
	e.ImagePath = m.ImagePath
	e.Email = m.Email
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt

	return e
}
