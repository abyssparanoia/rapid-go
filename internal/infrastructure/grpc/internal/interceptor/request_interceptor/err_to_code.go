package request_interceptor

import (
	"errors"

	pkg_errors "github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errToCode(err error) codes.Code {
	var notFoundError *pkg_errors.NotFoundError
	if ok := errors.As(err, &notFoundError); ok {
		return codes.NotFound
	}
	var invalidArgumentError *pkg_errors.BadRequestError
	if ok := errors.As(err, &invalidArgumentError); ok {
		return codes.InvalidArgument
	}
	var unauthenticatedError *pkg_errors.UnauthorizedError
	if ok := errors.As(err, &unauthenticatedError); ok {
		return codes.Unauthenticated
	}
	var alreadyExistError *pkg_errors.ConflictError
	if ok := errors.As(err, &alreadyExistError); ok {
		return codes.AlreadyExists
	}
	var internalError *pkg_errors.InternalError
	if ok := errors.As(err, &internalError); ok {
		return codes.Internal
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
	case codes.Internal, codes.Unknown, codes.Aborted:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
