package grpcerror

import (
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
