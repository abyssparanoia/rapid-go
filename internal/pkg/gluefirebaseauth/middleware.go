package gluefirebaseauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"

	"github.com/unrolled/render"
)

// Middleware ... http middleware for firebase authentication
type Middleware struct {
	firebaseauth Firebaseauth
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

		userID, claims, err := m.firebaseauth.Authentication(ctx, ah)
		if err != nil {
			m.renderError(ctx, w, http.StatusForbidden, err.Error())
			return
		}
		ctx = setUserID(ctx, userID)
		ctx = setClaims(ctx, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) renderError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	log.Warningf(ctx, msg)
	render.New().Text(w, status, fmt.Sprintf("%d authentication failed", status))
}

// NewMiddleware ... get middleware
func NewMiddleware(firebaseauth Firebaseauth) *Middleware {
	return &Middleware{
		firebaseauth,
	}
}
