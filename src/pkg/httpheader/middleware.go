package httpheader

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/pkg/log"

	"github.com/unrolled/render"
)

// Middleware ... middleware
type Middleware struct {
	Svc Service
}

// Handle ... get parameter from header and set to context
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p, err := m.Svc.Get(ctx, r)
		if err != nil {
			m.renderError(ctx, w, http.StatusBadRequest, "httpheader.Service.Get: "+err.Error())
			return
		}
		ctx = setParams(ctx, p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) renderError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	log.Warningf(ctx, msg)
	render.New().Text(w, status, fmt.Sprintf("%d invalid header params", status))
}

// NewMiddleware ... get middleware
func NewMiddleware(svc Service) *Middleware {
	return &Middleware{
		Svc: svc,
	}
}
