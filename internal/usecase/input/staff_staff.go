package input

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

// StaffGetStaff represents input for getting a staff member
type StaffGetStaff struct {
	StaffID       string    `validate:"required"`
	TenantID      string    `validate:"required"`
	TargetStaffID string    `validate:"required"`
	RequestTime   time.Time `validate:"required"`
}

func NewStaffGetStaff(
	staffID string,
	tenantID string,
	targetStaffID string,
	requestTime time.Time,
) *StaffGetStaff {
	return &StaffGetStaff{
		StaffID:       staffID,
		TenantID:      tenantID,
		TargetStaffID: targetStaffID,
		RequestTime:   requestTime,
	}
}

func (p *StaffGetStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

// StaffListStaffs represents input for listing staff members
type StaffListStaffs struct {
	StaffID     string `validate:"required"`
	TenantID    string `validate:"required"`
	Page        uint64
	Limit       uint64             `validate:"gte=1,lte=100"`
	SortKey     model.StaffSortKey // NON-nullable field
	RequestTime time.Time          `validate:"required"`
}

func NewStaffListStaffs(
	staffID string,
	tenantID string,
	page uint64,
	limit uint64,
	sortKey nullable.Type[model.StaffSortKey], // nullable param
	requestTime time.Time,
) *StaffListStaffs {
	// Pagination defaults
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 30
	}

	// SortKey default: CreatedAtDesc
	resolvedSortKey := model.StaffSortKeyCreatedAtDesc
	if sortKey.Valid && sortKey.Value().Valid() {
		resolvedSortKey = sortKey.Value()
	}

	return &StaffListStaffs{
		StaffID:     staffID,
		TenantID:    tenantID,
		Page:        page,
		Limit:       limit,
		SortKey:     resolvedSortKey,
		RequestTime: requestTime,
	}
}

func (p *StaffListStaffs) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
