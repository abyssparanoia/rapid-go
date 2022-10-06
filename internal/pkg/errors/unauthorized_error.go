package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// UnauthorizedError ... not found error
type UnauthorizedError string

func (e UnauthorizedError) Error() string {
	return string(e)
}

// New ... new error
func (e UnauthorizedError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf
func (e UnauthorizedError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error
func (e UnauthorizedError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf
func (e UnauthorizedError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// As ... as method
func (e UnauthorizedError) As(target interface{}) bool {
	if _, ok := target.(**UnauthorizedError); ok {
		return true
	}
	return false
}
