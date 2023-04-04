package input

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type PublicGetTenant struct {
	TenantID string `validate:"required"`
}

func NewPublicGetTenant(
	tenantID string,
) *PublicGetTenant {
	return &PublicGetTenant{
		TenantID: tenantID,
	}
}

func (p *PublicGetTenant) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
