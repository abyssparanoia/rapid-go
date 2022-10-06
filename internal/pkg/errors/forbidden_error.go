package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ForbiddenError ... forbidden error
type ForbiddenError string

func (e ForbiddenError) Error() string {
	return string(e)
}

// New ... new error
func (e ForbiddenError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf
func (e ForbiddenError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error
func (e ForbiddenError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf
func (e ForbiddenError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// As ... as method
func (e ForbiddenError) As(target interface{}) bool {
	if _, ok := target.(**ForbiddenError); ok {
		return true
	}
	return false
}
