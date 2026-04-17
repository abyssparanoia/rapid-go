package model

import (
	"time"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
)

type Admin struct {
	ID          string
	Role        AdminRole
	AuthUID     string
	Email       string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Admins []*Admin

func NewAdmin(
	role AdminRole,
	authUID string,
	email string,
	displayName string,
	t time.Time,
) *Admin {
	return &Admin{
		ID:          id.New(),
		Role:        role,
		AuthUID:     authUID,
		Email:       email,
		DisplayName: displayName,
		CreatedAt:   t,
		UpdatedAt:   t,
	}
}

func (m *Admin) Update(
	displayName null.String,
	t time.Time,
) {
	if displayName.Valid {
		m.DisplayName = displayName.String
	}

	m.UpdatedAt = t
}

func (m *Admin) UpdateRole(role AdminRole, t time.Time) *Admin {
	m.Role = role
	m.UpdatedAt = t
	return m
}

func (es Admins) IDs() []string {
	ids := make([]string, 0, len(es))
	for _, e := range es {
		ids = append(ids, e.ID)
	}
	return ids
}
