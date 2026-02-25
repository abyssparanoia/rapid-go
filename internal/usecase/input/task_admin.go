package input

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type TaskCreateAdmin struct {
	Email       string    `validate:"required,email"`
	DisplayName string    `validate:"required"`
	Password    string    `validate:"required"` //nolint:gosec
	RequestTime time.Time `validate:"required"`
}

func NewTaskCreateAdmin(
	email string,
	displayName string,
	password string,
	t time.Time,
) *TaskCreateAdmin {
	return &TaskCreateAdmin{
		Email:       email,
		DisplayName: displayName,
		Password:    password,
		RequestTime: t,
	}
}

func (p *TaskCreateAdmin) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	return nil
}
