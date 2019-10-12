package log

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go-worker/src/lib/util"
)

// Middleware ... middleware logger
type Middleware struct {
	Writer         Writer
	MinOutSeverity Severity
}

// Handle ... initialize logger
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startAt := util.TimeNow()

		traceID := util.StrUniqueID()
		logger := NewLogger(m.Writer, m.MinOutSeverity, traceID)
		ctx := r.Context()
		ctx = SetLogger(ctx, logger)

		defer func() {
			if rcvr := recover(); rcvr != nil {
				msg := Panic(ctx, rcvr)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(msg))

				endAt := util.TimeNow()
				dr := endAt.Sub(startAt)

				logger.WriteRequest(r, endAt, dr)
			}
		}()

		next.ServeHTTP(w, r.WithContext(ctx))

		endAt := util.TimeNow()
		dr := endAt.Sub(startAt)

		logger.WriteRequest(r, endAt, dr)
	})
}

// NewMiddleware ... get middleware
func NewMiddleware(writer Writer, minOutSeverity string) *Middleware {
	mos := NewSeverity(minOutSeverity)
	return &Middleware{
		Writer:         writer,
		MinOutSeverity: mos,
	}
}
