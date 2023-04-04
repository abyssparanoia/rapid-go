package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ServiceUnavailableError ... service unavailable error
type ServiceUnavailableError string

func (e ServiceUnavailableError) Error() string {
	return string(e)
}

// New ... new error
func (e ServiceUnavailableError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf
func (e ServiceUnavailableError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error
func (e ServiceUnavailableError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf
func (e ServiceUnavailableError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// As ... as method
func (e ServiceUnavailableError) As(target interface{}) bool {
	if _, ok := target.(**ServiceUnavailableError); ok {
		return true
	}
	return false
}
