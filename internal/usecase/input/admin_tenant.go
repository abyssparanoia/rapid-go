package input

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
	"github.com/volatiletech/null/v8"
)

type AdminGetTenant struct {
	TenantID string `validate:"required"`
}

func NewAdminGetTenant(
	tenantID string,
) *AdminGetTenant {
	return &AdminGetTenant{
		TenantID: tenantID,
	}
}

func (p *AdminGetTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type AdminListTenants struct {
	Page  uint64
	Limit uint64 `validate:"gte=1,lte=100"`
}

func NewAdminListTenants(
	page uint64,
	limit uint64,
) *AdminListTenants {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 30
	}
	return &AdminListTenants{
		Page:  page,
		Limit: limit,
	}
}

func (p *AdminListTenants) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type AdminCreateTenant struct {
	Name        string    `validate:"required"`
	RequestTime time.Time `validate:"required"`
}

func NewAdminCreateTenant(
	name string,
	requestTime time.Time,
) *AdminCreateTenant {
	return &AdminCreateTenant{
		Name:        name,
		RequestTime: requestTime,
	}
}

func (p *AdminCreateTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type AdminUpdateTenant struct {
	TenantID    string `validate:"required"`
	Name        null.String
	RequestTime time.Time `validate:"required"`
}

func NewAdminUpdateTenant(
	tenantID string,
	name null.String,
	requestTime time.Time,
) *AdminUpdateTenant {
	return &AdminUpdateTenant{
		TenantID:    tenantID,
		Name:        name,
		RequestTime: requestTime,
	}
}

func (p *AdminUpdateTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type AdminDeleteTenant struct {
	TenantID string `validate:"required"`
}

func NewAdminDeleteTenant(
	tenantID string,
) *AdminDeleteTenant {
	return &AdminDeleteTenant{
		TenantID: tenantID,
	}
}

func (p *AdminDeleteTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
