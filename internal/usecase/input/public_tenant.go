package input

import (
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/errors"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/validation"
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
