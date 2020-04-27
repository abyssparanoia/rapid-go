package httpmiddleware

import (
	"net/http"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/go-chi/chi/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HTTPMiddleware ... http middleware
type HTTPMiddleware struct {
	logger *zap.Logger
}

// Handle ... handle http request
func (m *HTTPMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var requestID string
		if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
			requestID = reqID.(string)
		}

		ctx := r.Context()
		ctx = ctxzap.ToContext(ctx, m.logger)

		defer func() {
			if rcvr := recover(); rcvr != nil {
			}
		}()

		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r.WithContext(ctx))

		latency := time.Since(start)

		fields := []zapcore.Field{
			zap.Int("status", sw.status),
			zap.Int("content-length", sw.length),
			zap.Duration("took", latency),
			zap.Int64("latency", latency.Nanoseconds()),
			zap.String("remote", r.RemoteAddr),
			zap.String("request", r.RequestURI),
			zap.String("method", r.Method),
		}
		if requestID != "" {
			fields = append(fields, zap.String("request-id", requestID))
		}
		log.Infof(ctx, "request completed", fields...)
	})
}

// New ... new http middleware
func New(logger *zap.Logger) *HTTPMiddleware {
	return &HTTPMiddleware{logger}
}
