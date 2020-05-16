package grpcerror

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error ... error
type Error struct {
	Code codes.Code
	Msg  string
}

func (e *Error) Error() string {
	return e.Msg
}

// New ... return wrapped error with blank string, which should `Wrap` pre-defined ErrXxx value (NOT errors.New).
func (e *Error) New() error {
	return status.New(e.Code, e.Msg).Err()
}

// Errorf ...
func (e *Error) Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return status.Errorf(e.Code, "err: %s; %s", e.Msg, msg)
}

// Wrap ...
func (e *Error) Wrap(err error) error {
	if err == nil {
		return e.New()
	}
	return status.Errorf(e.Code, "err: %s", e)
}

// Wrapf ...
func (e *Error) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return e.Errorf(format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return status.Errorf(e.Code, "err: %s; %s", err, msg)
}
