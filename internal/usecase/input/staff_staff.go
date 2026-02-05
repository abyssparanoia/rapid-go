package input

import (
	"time"

	"github.com/aarondl/null/v8"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

// StaffGetStaff represents input for getting a staff member
type StaffGetStaff struct {
	StaffID       string    `validate:"required"`
	TargetStaffID string    `validate:"required"`
	RequestTime   time.Time `validate:"required"`
}

func NewStaffGetStaff(
	staffID string,
	targetStaffID string,
	requestTime time.Time,
) *StaffGetStaff {
	return &StaffGetStaff{
		StaffID:       staffID,
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

type StaffCreateStaff struct {
	StaffID      string          `validate:"required"`
	TenantID     string          `validate:"required"`
	Email        string          `validate:"required"`
	DisplayName  string          `validate:"required"`
	Role         model.StaffRole `validate:"required"`
	ImageAssetID string          `validate:"required"`
	RequestTime  time.Time       `validate:"required"`
}

func NewStaffCreateStaff(
	staffID,
	tenantID,
	email,
	displayName string,
	role model.StaffRole,
	imageAssetID string,
	requestTime time.Time,
) *StaffCreateStaff {
	return &StaffCreateStaff{
		StaffID:      staffID,
		TenantID:     tenantID,
		Email:        email,
		DisplayName:  displayName,
		Role:         role,
		ImageAssetID: imageAssetID,
		RequestTime:  requestTime,
	}
}

func (p *StaffCreateStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

// StaffUpdateStaff represents input for updating a staff member
type StaffUpdateStaff struct {
	StaffID       string `validate:"required"`
	TargetStaffID string `validate:"required"`
	DisplayName   null.String
	Role          nullable.Type[model.StaffRole]
	ImageAssetID  null.String
	RequestTime   time.Time `validate:"required"`
}

func NewStaffUpdateStaff(
	staffID string,
	targetStaffID string,
	displayName null.String,
	role nullable.Type[model.StaffRole],
	imageAssetID null.String,
	requestTime time.Time,
) *StaffUpdateStaff {
	return &StaffUpdateStaff{
		StaffID:       staffID,
		TargetStaffID: targetStaffID,
		DisplayName:   displayName,
		Role:          role,
		ImageAssetID:  imageAssetID,
		RequestTime:   requestTime,
	}
}

func (p *StaffUpdateStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	if p.Role.Valid && !p.Role.Value().Valid() {
		return errors.RequestInvalidArgumentErr.Errorf("invalid role: %s", p.Role.Value())
	}
	return nil
}
