package log

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/src/config"
	"go.uber.org/zap"
)

type loggerKey struct{}

// NewLogger ... create logger and set it in contenxt
func NewLogger(ctx context.Context) context.Context {

	// development
	if config.IsEnvDeveloping() {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		ctx = context.WithValue(ctx, loggerKey{}, logger)
		// production
	} else {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		ctx = context.WithValue(ctx, loggerKey{}, logger)
	}
	return ctx
}

// Logger ... get context from context
func Logger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(loggerKey{}).(*zap.Logger)

	return logger
}

// TODO: change appropriate using

// Debugf ... output debug log
func Debugf(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Debug(fmt.Sprintf(msg, fields...))
}

// Infof ... output info log
func Infof(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Info(fmt.Sprintf(msg, fields...))
}

// Warningf ... output warning log
func Warningf(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Warn(fmt.Sprintf(msg, fields...))
}

// Errorf ... output error log
func Errorf(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Error(fmt.Sprintf(msg, fields...))
}
