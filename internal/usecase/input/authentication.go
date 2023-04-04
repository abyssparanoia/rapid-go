package input

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type VerifyIDToken struct {
	IDToken string `validate:"required"`
}

func NewVerifyIDToken(
	idToken string,
) *VerifyIDToken {
	return &VerifyIDToken{
		IDToken: idToken,
	}
}

func (p *VerifyIDToken) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
