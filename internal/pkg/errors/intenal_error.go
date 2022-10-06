package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// InternalError ... base error
type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

// New ... new error
func (e InternalError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf
func (e InternalError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error
func (e InternalError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf
func (e InternalError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}
