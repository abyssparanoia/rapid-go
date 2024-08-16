package input

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type CreateRootStaff struct {
	Email       string    `validate:"required"`
	Password    string    `validate:"required"`
	RequestTime time.Time `validate:"required"`
}

func NewCreateRootStaff(
	email string,
	password string,
	t time.Time,
) *CreateRootStaff {
	return &CreateRootStaff{
		Email:       email,
		Password:    password,
		RequestTime: t,
	}
}

func (p *CreateRootStaff) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
