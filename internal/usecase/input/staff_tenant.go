package input

import (
	"time"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type StaffGetTenant struct {
	StaffID  string `validate:"required"`
	TenantID string `validate:"required"`
}

func NewStaffGetTenant(
	staffID string,
	tenantID string,
) *StaffGetTenant {
	return &StaffGetTenant{
		StaffID:  staffID,
		TenantID: tenantID,
	}
}

func (p *StaffGetTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type StaffListTenants struct {
	StaffID string
	Page    uint64
	Limit   uint64              `validate:"gte=1,lte=100"`
	SortKey model.TenantSortKey // NON-nullable field
}

func NewStaffListTenants(
	staffID string,
	page uint64,
	limit uint64,
	sortKey nullable.Type[model.TenantSortKey], // nullable param
) *StaffListTenants {
	// Pagination defaults
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 30
	}

	// SortKey default: CreatedAtDesc
	resolvedSortKey := model.TenantSortKeyCreatedAtDesc
	if sortKey.Valid && sortKey.Value().Valid() {
		resolvedSortKey = sortKey.Value()
	}

	return &StaffListTenants{
		StaffID: staffID,
		Page:    page,
		Limit:   limit,
		SortKey: resolvedSortKey,
	}
}

func (p *StaffListTenants) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type StaffCreateTenant struct {
	StaffID     string    `validate:"required"`
	Name        string    `validate:"required"`
	RequestTime time.Time `validate:"required"`
}

func NewStaffCreateTenant(
	staffID string,
	name string,
	requestTime time.Time,
) *StaffCreateTenant {
	return &StaffCreateTenant{
		StaffID:     staffID,
		Name:        name,
		RequestTime: requestTime,
	}
}

func (p *StaffCreateTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type StaffUpdateTenant struct {
	StaffID     string `validate:"required"`
	TenantID    string `validate:"required"`
	Name        null.String
	RequestTime time.Time `validate:"required"`
}

func NewStaffUpdateTenant(
	staffID string,
	tenantID string,
	name null.String,
	requestTime time.Time,
) *StaffUpdateTenant {
	return &StaffUpdateTenant{
		StaffID:     staffID,
		TenantID:    tenantID,
		Name:        name,
		RequestTime: requestTime,
	}
}

func (p *StaffUpdateTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type StaffDeleteTenant struct {
	StaffID  string `validate:"required"`
	TenantID string `validate:"required"`
}

func NewStaffDeleteTenant(
	staffID string,
	tenantID string,
) *StaffDeleteTenant {
	return &StaffDeleteTenant{
		StaffID:  staffID,
		TenantID: tenantID,
	}
}

func (p *StaffDeleteTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
