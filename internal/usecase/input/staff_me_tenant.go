package input

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type StaffGetMeTenant struct {
	TenantID    string    `validate:"required"`
	StaffID     string    `validate:"required"`
	RequestTime time.Time `validate:"required"`
}

func NewStaffGetMeTenant(
	tenantID string,
	staffID string,
	requestTime time.Time,
) *StaffGetMeTenant {
	return &StaffGetMeTenant{
		TenantID:    tenantID,
		StaffID:     staffID,
		RequestTime: requestTime,
	}
}

func (p *StaffGetMeTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type StaffUpdateMeTenant struct {
	TenantID    string `validate:"required"`
	StaffID     string `validate:"required"`
	Name        null.String
	RequestTime time.Time `validate:"required"`
}

func NewStaffUpdateMeTenant(
	tenantID string,
	staffID string,
	name null.String,
	requestTime time.Time,
) *StaffUpdateMeTenant {
	return &StaffUpdateMeTenant{
		TenantID:    tenantID,
		StaffID:     staffID,
		Name:        name,
		RequestTime: requestTime,
	}
}

func (p *StaffUpdateMeTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
