package input

import (
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/errors"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/validation"
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
