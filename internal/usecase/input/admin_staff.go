package input

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type AdminCreateStaff struct {
	TenantID    string          `validate:"required"`
	Email       string          `validate:"required"`
	DisplayName string          `validate:"required"`
	Role        model.StaffRole `validate:"required"`
	AssetKey    string          `validate:"required"`
	RequestTime time.Time       `validate:"required"`
}

func NewAdminCreateStaff(
	tenantID,
	email,
	displayName string,
	role model.StaffRole,
	assetKey string,
	requestTime time.Time,
) *AdminCreateStaff {
	return &AdminCreateStaff{
		TenantID:    tenantID,
		Email:       email,
		DisplayName: displayName,
		Role:        role,
		AssetKey:    assetKey,
		RequestTime: requestTime,
	}
}

func (p *AdminCreateStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
