package log

import "context"

type contextKey string

type loggerContextKey struct{}

// GetLogger ... get logger from context
func GetLogger(ctx context.Context) *Logger {
	if itf := ctx.Value(loggerContextKey{}); itf != nil {
		logger := itf.(*Logger)
		return logger
	}
	return nil
}

// SetLogger ... set logger to context
func SetLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}
