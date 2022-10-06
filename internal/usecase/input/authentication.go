package input

import (
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/errors"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/validation"
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
