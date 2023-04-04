package input

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type PublicSignIn struct {
	AuthUID string `validate:"required"`
}

func NewPublicSignIn(
	authUID string,
) *PublicSignIn {
	return &PublicSignIn{
		AuthUID: authUID,
	}
}

func (p *PublicSignIn) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
