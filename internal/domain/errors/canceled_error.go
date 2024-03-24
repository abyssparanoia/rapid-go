package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// CanceledError ... base error.
type CanceledError string

func (e CanceledError) Error() string {
	return string(e)
}

// New ... new error.
func (e CanceledError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf.
func (e CanceledError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error.
func (e CanceledError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf.
func (e CanceledError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// As ... as method.
func (e CanceledError) As(target interface{}) bool {
	if _, ok := target.(**CanceledError); ok {
		return true
	}
	return false
}
