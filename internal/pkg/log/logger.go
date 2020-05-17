package log

import (
	"context"

	"go.uber.org/zap/zapcore"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

// Must ... new logger
func Must(ctx context.Context) *zap.Logger {
	return logger(ctx)
}

// New ... new logger
func New(env string) (*zap.Logger, error) {
	return newStackdriverConfig().Build(zap.AddStacktrace(zapcore.WarnLevel))
}

// Logger ... get context from context
func logger(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}

// Debugf ... output debug log
func Debugf(ctx context.Context, msg string, fields ...zap.Field) {
	withTracing(
		ctx,
		msg,
		func(ctx context.Context, msg string, fields ...zap.Field) {
			logger(ctx).WithOptions(zap.AddCallerSkip(3)).Debug(msg, fields...)
		},
		fields...)
}

// SugarDebugf ... output sugar debug log
func SugarDebugf(ctx context.Context, msg string, args ...interface{}) {
	logger(ctx).WithOptions(zap.AddCallerSkip(3)).Sugar().Debugf(msg, args...)
}

// Infof ... output info log
func Infof(ctx context.Context, msg string, fields ...zap.Field) {
	withTracing(
		ctx,
		msg,
		func(ctx context.Context, msg string, fields ...zap.Field) {
			logger(ctx).WithOptions(zap.AddCallerSkip(3)).Info(msg, fields...)
		},
		fields...)
}

// Warningf ... output warning log
func Warningf(ctx context.Context, msg string, fields ...zap.Field) {
	withTracing(
		ctx,
		msg,
		func(ctx context.Context, msg string, fields ...zap.Field) {
			logger(ctx).WithOptions(zap.AddCallerSkip(3)).Warn(msg, fields...)
		},
		fields...)
}

// Errorf ... output error log
func Errorf(ctx context.Context, msg string, fields ...zap.Field) {
	withTracing(
		ctx,
		msg,
		func(ctx context.Context, msg string, fields ...zap.Field) {
			logger(ctx).WithOptions(zap.AddCallerSkip(3)).Error(msg, fields...)
		},
		fields...)
}

// ErrorfIfExists ... calls Errorf only when the error exists
func ErrorfIfExists(ctx context.Context, err error, msg string, fields ...zap.Field) {
	if err == nil {
		return
	}
	withTracing(
		ctx,
		msg,
		func(ctx context.Context, msg string, fields ...zap.Field) {
			logger(ctx).WithOptions(zap.AddCallerSkip(3)).Error(msg, fields...)
		},
		fields...)
}

// LogFunc ... log func type
type LogFunc func(ctx context.Context, msg string, fields ...zap.Field)

func withTracing(
	ctx context.Context,
	msg string,
	f LogFunc,
	fields ...zap.Field,
) {
	if len(fields) == 0 {
		fields = make([]zap.Field, 0)
	}
	if ctx != nil {
		sp := opentracing.SpanFromContext(ctx)
		if sp != nil {
			spc := sp.Context().(jaeger.SpanContext)

			fields = append(fields, zap.String("TraceID", spc.TraceID().String()))
			fields = append(fields, zap.String("ParentID", spc.ParentID().String()))
			fields = append(fields, zap.String("SpanID", spc.SpanID().String()))
		}
	}
	f(ctx, msg, fields...)
}
