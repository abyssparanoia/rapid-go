package logger

import (
	"context"

	"github.com/blendle/zapdriver"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

// New ... new logger.
func New() *zap.Logger {
	l, err := zapdriver.NewProduction()
	if err != nil {
		panic(err)
	}
	return l
}

// L ... logger.
var L = func(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}

func AddFields(ctx context.Context, fields ...zap.Field) {
	ctxzap.AddFields(ctx, fields...)
}

func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return ctxzap.ToContext(ctx, logger)
}
