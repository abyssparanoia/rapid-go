package model

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/volatiletech/null/v8"
)

type Tenant struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tenants []*Tenant

func NewTenant(
	name string,
	t time.Time,
) *Tenant {
	return &Tenant{
		ID:        id.New(),
		Name:      name,
		CreatedAt: t,
		UpdatedAt: t,
	}
}

func (m *Tenant) Update(
	name null.String,
	t time.Time,
) {
	if name.Valid {
		m.Name = name.String
	}

	m.UpdatedAt = t
}
