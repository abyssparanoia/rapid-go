package input

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type CreateRootUser struct {
	Email       string    `validate:"required"`
	Passoword   string    `validate:"required"`
	RequestTime time.Time `validate:"required"`
}

func NewCreateRootUser(
	email string,
	password string,
	t time.Time,
) *CreateRootUser {
	return &CreateRootUser{
		Email:       email,
		Passoword:   password,
		RequestTime: t,
	}
}

func (p *CreateRootUser) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
