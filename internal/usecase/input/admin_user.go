package input

import (
	"time"

	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/errors"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/validation"
)

type AdminCreateUser struct {
	TenantID    string         `validate:"required"`
	Email       string         `validate:"required"`
	DisplayName string         `validate:"required"`
	Role        model.UserRole `validate:"required"`
	RequestTime time.Time      `validate:"required"`
}

func NewAdminCreateUser(
	tenantID,
	email,
	displayName string,
	role model.UserRole,
	requestTime time.Time,
) *AdminCreateUser {
	return &AdminCreateUser{
		TenantID:    tenantID,
		Email:       email,
		DisplayName: displayName,
		Role:        role,
		RequestTime: requestTime,
	}
}

func (p *AdminCreateUser) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
