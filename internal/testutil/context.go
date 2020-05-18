package testutil

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

var (
	testCtx context.Context
)

func init() {
	logger, _ := log.New("test")
	testCtx = ctxzap.ToContext(context.Background(), logger)
}

// Context ... context for test
func Context() context.Context {
	return testCtx
}
