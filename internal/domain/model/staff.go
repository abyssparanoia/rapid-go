package model

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/ulid"
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

	Tenant *Tenant
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
		ID:          ulid.New(),
		TenantID:    tenantID,
		Role:        role,
		AuthUID:     authUID,
		DisplayName: displayName,
		ImagePath:   imagePath,
		Email:       email,
		CreatedAt:   t,
		UpdatedAt:   t,
	}
}

func (m *Staff) Exist() bool {
	return m != nil
}
