package model

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
)

type Tenant struct {
	ID        string
	Name      string
	Tags      TenantTags
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
		Tags:      make(TenantTags, 0),
		CreatedAt: t,
		UpdatedAt: t,
	}
}

type TenantTag struct {
	ID        string
	Type      TenantTagType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TenantTags []*TenantTag

func NewTenantTag(
	tagType TenantTagType,
	t time.Time,
) *TenantTag {
	return &TenantTag{
		ID:        id.New(),
		Type:      tagType,
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
