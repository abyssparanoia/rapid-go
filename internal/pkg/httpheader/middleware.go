package httpheader

import (
	"context"
	"fmt"
	"net/http"

	"github.com/unrolled/render"
)

// Middleware ... middleware
type Middleware struct {
	httpheader Httpheader
}

// Handle ... get parameter from header and set to context
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p, err := m.httpheader.Get(ctx, r)
		if err != nil {
			m.renderError(ctx, w, http.StatusBadRequest, "httpheader.Service.Get: "+err.Error())
			return
		}
		ctx = setParams(ctx, p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) renderError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	render.New().Text(w, status, fmt.Sprintf("%d invalid header params", status))
}

// NewMiddleware ... get middleware
func NewMiddleware(httpheader Httpheader) *Middleware {
	return &Middleware{
		httpheader,
	}
}
