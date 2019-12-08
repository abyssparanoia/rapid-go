package firebaseauth

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

// Handle ... authenticate handler
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ah := r.Header.Get("Authorization")
		if ah == "" {
			m.renderError(ctx, w, http.StatusForbidden, "no Authorization header")
			return
		}
		ctx = setAuthHeader(ctx, ah)

		userID, claims, err := m.Svc.Authentication(ctx, ah)
		if err != nil {
			m.renderError(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		ctx = setUserID(ctx, userID)
		log.Debugf(ctx, "UserID: %s", userID)

		ctx = setClaims(ctx, claims)
		log.Debugf(ctx, "Claims: %v", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) renderError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	log.Warningf(ctx, msg)
	render.New().Text(w, status, fmt.Sprintf("%d authentication failed", status))
}

// NewMiddleware ... get middleware
func NewMiddleware(svc Service) *Middleware {
	return &Middleware{
		Svc: svc,
	}
}
