package logger

import (
	"context"

	"github.com/blendle/zapdriver"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

// New ... new logger
func New() *zap.Logger {
	l, err := zapdriver.NewProduction()
	if err != nil {
		panic(err)
	}
	return l
}

// L ... logger
var L = func(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}
