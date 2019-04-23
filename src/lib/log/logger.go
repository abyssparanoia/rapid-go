package log

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/config"
	"go.uber.org/zap"
)

type loggerKey struct{}

// NewLogger ... create logger and set it in contenxt in first middleware
func NewLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger ... get context from context
func logger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(loggerKey{}).(*zap.Logger)

	// for test code
	if logger == nil {
		logger, _ := zap.NewDevelopment()
		return logger
	}

	return logger
}

// TODO: change appropriate using

// Debugf ... output debug log
func Debugf(ctx context.Context, msg string, fields ...interface{}) {
	logger(ctx).Debug(fmt.Sprintf(msg, fields...))
}

// Infof ... output info log
func Infof(ctx context.Context, msg string, fields ...interface{}) {
	logger(ctx).Info(fmt.Sprintf(msg, fields...))
}

// Warningf ... output warning log
func Warningf(ctx context.Context, msg string, fields ...interface{}) {
	logger(ctx).Warn(fmt.Sprintf(msg, fields...))
}

// Errorf ... output error log
func Errorf(ctx context.Context, msg string, fields ...interface{}) {
	logger(ctx).Error(fmt.Sprintf(msg, fields...))
}
