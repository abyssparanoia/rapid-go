package input

import (
	"time"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type StaffSignUp struct {
	AuthUID      string    `validate:"required"`
	Email        string    `validate:"required,email"`
	TenantName   string    `validate:"required,max=256"`
	DisplayName  string    `validate:"required,max=256"`
	ImageAssetID string    `validate:"required"`
	RequestTime  time.Time `validate:"required"`
}

func NewStaffSignUp(
	authUID string,
	email string,
	tenantName string,
	displayName string,
	imageAssetID string,
	requestTime time.Time,
) *StaffSignUp {
	return &StaffSignUp{
		AuthUID:      authUID,
		Email:        email,
		TenantName:   tenantName,
		DisplayName:  displayName,
		ImageAssetID: imageAssetID,
		RequestTime:  requestTime,
	}
}

func (p *StaffSignUp) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type StaffGetMe struct {
	TenantID    string    `validate:"required"`
	StaffID     string    `validate:"required"`
	RequestTime time.Time `validate:"required"`
}

func NewStaffGetMe(
	tenantID string,
	staffID string,
	requestTime time.Time,
) *StaffGetMe {
	return &StaffGetMe{
		TenantID:    tenantID,
		StaffID:     staffID,
		RequestTime: requestTime,
	}
}

func (p *StaffGetMe) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type StaffUpdateMe struct {
	TenantID     string `validate:"required"`
	StaffID      string `validate:"required"`
	DisplayName  null.String
	ImageAssetID null.String
	RequestTime  time.Time `validate:"required"`
}

func NewStaffUpdateMe(
	tenantID string,
	staffID string,
	displayName null.String,
	imageAssetID null.String,
	requestTime time.Time,
) *StaffUpdateMe {
	return &StaffUpdateMe{
		TenantID:     tenantID,
		StaffID:      staffID,
		DisplayName:  displayName,
		ImageAssetID: imageAssetID,
		RequestTime:  requestTime,
	}
}

func (p *StaffUpdateMe) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
