package model

import (
	"time"

	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/ulid"
)

type User struct {
	ID          string
	Tenant      *Tenant
	Role        UserRole
	AuthUID     string
	DisplayName string
	ImagePath   string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUser(
	tenantID string,
	role UserRole,
	authUID string,
	displayName string,
	imagePath string,
	email string,
	t time.Time,
) *User {
	return &User{
		ID: ulid.New(),
		Tenant: &Tenant{
			ID: tenantID,
		},
		Role:        role,
		AuthUID:     authUID,
		DisplayName: displayName,
		ImagePath:   imagePath,
		Email:       email,
		CreatedAt:   t,
		UpdatedAt:   t,
	}
}

func (m *User) Exist() bool {
	return m != nil
}
