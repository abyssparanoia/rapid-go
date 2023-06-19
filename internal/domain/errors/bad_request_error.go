package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// BadRequestError ... bad request error.
type BadRequestError string

func (e BadRequestError) Error() string {
	return string(e)
}

// New ... new error.
func (e BadRequestError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf.
func (e BadRequestError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error.
func (e BadRequestError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf.
func (e BadRequestError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// As ... as method.
func (e BadRequestError) As(target interface{}) bool {
	if _, ok := target.(**BadRequestError); ok {
		return true
	}
	return false
}
