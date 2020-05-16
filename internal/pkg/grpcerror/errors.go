package grpcerror

import (
	"google.golang.org/grpc/codes"
)

// ErrUserNotFound ... err user not found
func ErrUserNotFound() *Error {
	return &Error{
		Code: codes.NotFound,
		Msg:  "ErrUserNotFound",
	}
}
