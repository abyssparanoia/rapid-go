package httperror

import (
	"fmt"

	"github.com/pkg/errors"
)

// HTTPError ... base error
type HTTPError string

func (e HTTPError) Error() string {
	return string(e)
}

// New ... new error
func (e HTTPError) New() error {
	return errors.Wrap(e, "")
}

// Errorf ... errorf
func (e HTTPError) Errorf(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// Wrap ... wrap error
func (e HTTPError) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return errors.Wrap(e, err.Error())
}

// Wrapf ... wrapf
func (e HTTPError) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.Wrapf(e, "err: %s; %s", err, msg)
}

// NotFoundError ... not found error
type NotFoundError = HTTPError

// As ... as method
func (ne NotFoundError) As(target interface{}) bool {
	if _, ok := target.(**NotFoundError); ok {
		return true
	}
	return false
}
