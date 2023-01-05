package errors

import "github.com/pkg/errors"

const (
	InternalErr               InternalError     = "An internal error has occurred2"
	UnauthorizedErr           UnauthorizedError = "Unauthroized"
	RequestInvalidArgumentErr BadRequestError   = "Request argument is invalid"
	NotFoundErr               NotFoundError     = "Not found"
)

func ExtractPlaneErrMessage(err error) string {
	switch errors.Cause(err) {
	case UnauthorizedErr:
		return UnauthorizedErr.Error()
	case RequestInvalidArgumentErr:
		return RequestInvalidArgumentErr.Error()
	case NotFoundErr:
		return NotFoundErr.Error()
	}
	return InternalErr.Error()
}
