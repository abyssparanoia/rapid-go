package logger

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxMarker struct{}

type ctxLogger struct {
	logger *zap.Logger
	fields []zapcore.Field
}

var (
	ctxMarkerKey = &ctxMarker{}
	nullLogger   = zap.NewNop()
)

// New ... new logger.
func New(env environment.MinLogLevel) *zap.Logger {
	// Set log level based on environment
	var level zapcore.Level
	switch env {
	case environment.MinLogLevelDebug:
		level = zapcore.DebugLevel
	case environment.MinLogLevelInfo:
		level = zapcore.InfoLevel
	case environment.MinLogLevelWarning:
		level = zapcore.WarnLevel
	default:
		level = zapcore.InfoLevel
	}

	config := zapdriver.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)

	l, err := config.Build(zapdriver.WrapCore())
	if err != nil {
		panic(err)
	}
	return l
}

// L ... logger.
var L = Extract

func AddFields(ctx context.Context, fields ...zapcore.Field) {
	l, ok := ctx.Value(ctxMarkerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	l.fields = append(l.fields, fields...)
}

func Extract(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(ctxMarkerKey).(*ctxLogger)
	if !ok || l == nil {
		return nullLogger
	}
	fields := TagsToFields(ctx)
	fields = append(fields, l.fields...)
	return l.logger.With(fields...)
}

func TagsToFields(ctx context.Context) []zapcore.Field {
	fields := []zapcore.Field{}
	return fields
}

func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	l := &ctxLogger{
		logger: logger,
	}
	return context.WithValue(ctx, ctxMarkerKey, l)
}
