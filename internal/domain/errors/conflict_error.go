package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ConflictError ... business error.
type ConflictError string

func (e ConflictError) Error() string {
	return string(e)
}

// New ... new error.
func (e ConflictError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf.
func (e ConflictError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error.
func (e ConflictError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf.
func (e ConflictError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// As ... as method.
func (e ConflictError) As(target interface{}) bool {
	if _, ok := target.(**ConflictError); ok {
		return true
	}
	return false
}
