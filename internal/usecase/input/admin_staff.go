package input

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/nullable"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

// AdminGetStaff represents input for getting a staff member
type AdminGetStaff struct {
	StaffID     string    `validate:"required"`
	RequestTime time.Time `validate:"required"`
}

func NewAdminGetStaff(
	staffID string,
	requestTime time.Time,
) *AdminGetStaff {
	return &AdminGetStaff{
		StaffID:     staffID,
		RequestTime: requestTime,
	}
}

func (p *AdminGetStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

// AdminListStaffs represents input for listing staff members
type AdminListStaffs struct {
	TenantID    string `validate:"required"`
	Page        uint64
	Limit       uint64             `validate:"gte=1,lte=100"`
	SortKey     model.StaffSortKey // NON-nullable field
	RequestTime time.Time          `validate:"required"`
}

func NewAdminListStaffs(
	tenantID string,
	page uint64,
	limit uint64,
	sortKey nullable.Type[model.StaffSortKey], // nullable param
	requestTime time.Time,
) *AdminListStaffs {
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

	return &AdminListStaffs{
		TenantID:    tenantID,
		Page:        page,
		Limit:       limit,
		SortKey:     resolvedSortKey,
		RequestTime: requestTime,
	}
}

func (p *AdminListStaffs) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

type AdminCreateStaff struct {
	AdminID      string          `validate:"required"`
	TenantID     string          `validate:"required"`
	Email        string          `validate:"required"`
	DisplayName  string          `validate:"required"`
	Role         model.StaffRole `validate:"required"`
	ImageAssetID string          `validate:"required"`
	RequestTime  time.Time       `validate:"required"`
}

func NewAdminCreateStaff(
	adminID,
	tenantID,
	email,
	displayName string,
	role model.StaffRole,
	imageAssetID string,
	requestTime time.Time,
) *AdminCreateStaff {
	return &AdminCreateStaff{
		AdminID:      adminID,
		TenantID:     tenantID,
		Email:        email,
		DisplayName:  displayName,
		Role:         role,
		ImageAssetID: imageAssetID,
		RequestTime:  requestTime,
	}
}

func (p *AdminCreateStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}

// AdminUpdateStaff represents input for updating a staff member
type AdminUpdateStaff struct {
	AdminID      string `validate:"required"`
	StaffID      string `validate:"required"`
	DisplayName  null.String
	Role         nullable.Type[model.StaffRole]
	ImageAssetID null.String
	RequestTime  time.Time `validate:"required"`
}

func NewAdminUpdateStaff(
	adminID string,
	staffID string,
	displayName null.String,
	role nullable.Type[model.StaffRole],
	imageAssetID null.String,
	requestTime time.Time,
) *AdminUpdateStaff {
	return &AdminUpdateStaff{
		AdminID:      adminID,
		StaffID:      staffID,
		DisplayName:  displayName,
		Role:         role,
		ImageAssetID: imageAssetID,
		RequestTime:  requestTime,
	}
}

func (p *AdminUpdateStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	if p.Role.Valid && !p.Role.Value().Valid() {
		return errors.RequestInvalidArgumentErr.Errorf("invalid role: %s", p.Role.Value())
	}
	return nil
}
