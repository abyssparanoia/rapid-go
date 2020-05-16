package grpcerror

import (
	"errors"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

// CodeToLevel ... convert grpc code to zapcore level
func CodeToLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.NotFound, codes.InvalidArgument, codes.AlreadyExists, codes.Unauthenticated, codes.PermissionDenied:
		return zapcore.WarnLevel
	case codes.Internal, codes.Unknown, codes.Aborted:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ExtractCodeFromErr ... extract code from error
func ExtractCodeFromErr(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	var se *Error
	if ok := errors.Is(err, se); ok {
		return se.Code
	}
	return codes.Unknown
}
