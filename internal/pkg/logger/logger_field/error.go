package logger_field

import (
	"github.com/abyssparanoia/goerr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Error(err error) zapcore.Field {
	if goErr := goerr.Unwrap(err); goErr != nil {
		return zap.Object("error", goErr)
	}
	return zap.Error(err) //nolint:forbidigo
}
