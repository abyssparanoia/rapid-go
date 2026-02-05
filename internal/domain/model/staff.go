package model

import (
	"time"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
)

type Staff struct {
	ID          string
	TenantID    string
	Role        StaffRole
	AuthUID     string
	DisplayName string
	ImagePath   string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	ReadonlyReference *struct {
		Tenant *Tenant
	}

	ImageURL null.String
}

type Staffs []*Staff

func NewStaff(
	tenantID string,
	role StaffRole,
	authUID string,
	displayName string,
	imagePath string,
	email string,
	t time.Time,
) *Staff {
	return &Staff{
		ID:          id.New(),
		TenantID:    tenantID,
		Role:        role,
		AuthUID:     authUID,
		DisplayName: displayName,
		ImagePath:   imagePath,
		Email:       email,
		CreatedAt:   t,
		UpdatedAt:   t,

		ReadonlyReference: nil,

		ImageURL: null.String{},
	}
}

func (m *Staff) Exist() bool {
	return m != nil
}

func (m *Staff) SetImageURL(
	imageURL string,
) {
	m.ImageURL = null.StringFrom(imageURL)
}

func (m *Staff) Update(
	displayName null.String,
	role nullable.Type[StaffRole],
	imagePath null.String,
	t time.Time,
) {
	if displayName.Valid {
		m.DisplayName = displayName.String
	}
	if role.Valid {
		m.Role = role.Value()
	}
	if imagePath.Valid {
		m.ImagePath = imagePath.String
	}
	m.UpdatedAt = t
}
