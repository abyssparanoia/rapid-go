package httpheader

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/gke-beego/api/src/lib/log"

	"github.com/unrolled/render"
)

// Middleware ... Headerに関する機能を提供する
type Middleware struct {
	Svc Service
}

// Handle ... リクエストヘッダーのパラメータを取得する
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

// NewMiddleware ... Middlewareを作成する
func NewMiddleware(svc Service) *Middleware {
	return &Middleware{
		Svc: svc,
	}
}
