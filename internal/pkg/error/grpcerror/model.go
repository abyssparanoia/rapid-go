package grpcerror

import (
	"fmt"

	"github.com/pkg/errors"
)

// Error ... base error
type Error string

func (e Error) Error() string {
	return string(e)
}

// New ... new error
func (e Error) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf
func (e Error) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error
func (e Error) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf
func (e Error) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// NotFoundError ... not found error
type NotFoundError = Error

// As ... as method
func (ne NotFoundError) As(target interface{}) bool {
	if _, ok := target.(**NotFoundError); ok {
		return true
	}
	return false
}
