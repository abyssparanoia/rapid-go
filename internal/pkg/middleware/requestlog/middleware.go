package requestlog

import (
	"net/http"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/error/httperror"

	"github.com/blendle/zapdriver"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const producerID = "rapid-go"

// HTTPMiddleware ... http middleware
type HTTPMiddleware struct {
	logger *zap.Logger
}

// Handle ... handle http request
func (m *HTTPMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		operationID := uuid.New()

		m.logger.Info("call start", zapdriver.OperationStart(operationID.String(), producerID))

		ctx := r.Context()
		ctx = ctxzap.ToContext(ctx, m.logger.With(
			zapdriver.OperationCont(operationID.String(), producerID),
		))

		defer func() {
			if rcvr := recover(); rcvr != nil {
			}
		}()

		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r.WithContext(ctx))

		latency := time.Since(start)
		zapcoreLevel := httperror.CodeToLevel(sw.status)
		ctxzap.Extract(ctx).Check(zapcoreLevel, "call end").Write(
			zap.Int("status", sw.status),
			zap.Int("content-length", sw.length),
			zap.Duration("took", latency),
			zap.Int64("latency", latency.Nanoseconds()),
			zap.String("remote", r.RemoteAddr),
			zap.String("request", r.RequestURI),
			zap.String("method", r.Method),
			zapdriver.OperationEnd(operationID.String(), producerID))
	})
}

// New ... new http middleware
func New(logger *zap.Logger) *HTTPMiddleware {
	return &HTTPMiddleware{logger}
}
