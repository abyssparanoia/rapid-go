package request_interceptor

import (
	"github.com/abyssparanoia/goerr"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errToCode(err error) codes.Code {
	if goErr := goerr.Unwrap(err); goErr != nil {
		switch goErr.Category() {
		case errors.ErrorCategoryNotFound.String():
			return codes.NotFound
		case errors.ErrorCategoryBadRequest.String():
			return codes.InvalidArgument
		case errors.ErrorCategoryUnauthorized.String():
			return codes.Unauthenticated
		case errors.ErrorCategoryConflict.String():
			return codes.AlreadyExists
		case errors.ErrorCategoryInternal.String():
			return codes.Internal
		}
	}
	if status, ok := status.FromError(err); ok {
		return status.Code()
	}
	return codes.Unknown
}

func codeToZapCoreLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.NotFound, codes.InvalidArgument, codes.AlreadyExists, codes.Unauthenticated, codes.PermissionDenied, codes.FailedPrecondition:
		return zapcore.WarnLevel
	case codes.Internal, codes.Unknown, codes.Aborted, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Unavailable, codes.Canceled, codes.OutOfRange, codes.Unimplemented, codes.DataLoss:
		return zapcore.ErrorLevel
	case codes.OK:
		fallthrough
	default:
		return zapcore.InfoLevel
	}
}
