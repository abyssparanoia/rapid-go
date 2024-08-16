package model

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/id"
	"github.com/volatiletech/null/v8"
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

	ImageURL null.String
	Tenant   *Tenant
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

		ImageURL: null.String{},
		Tenant:   nil,
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
