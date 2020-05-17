package httperror

import (
	"errors"
	"net/http"

	"go.uber.org/zap/zapcore"
)

// CodeToLevel ... convert grpc code to zapcore level
func CodeToLevel(code int) zapcore.Level {
	switch code {
	case http.StatusNotFound, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		return zapcore.WarnLevel
	case http.StatusInternalServerError, http.StatusTooManyRequests, http.StatusServiceUnavailable:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ErrToCode ... convert err to code
func ErrToCode(err error) int {
	var notFoundError *NotFoundError
	if ok := errors.As(err, &notFoundError); ok {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
