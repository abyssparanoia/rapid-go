package basicauth

import (
	"net/http"
)

// Middleware ... middleware for checking basic auth
type Middleware struct {
	Account *Account
}

// Handle ... handler
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "basic auth required.", http.StatusUnauthorized)
			return
		}
		if m.Account.User == user && m.Account.Password == password {
			http.Error(w, "basic auth error.", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NewMiddleware ... get middleware
func NewMiddleware(account *Account) *Middleware {
	return &Middleware{
		Account: account,
	}
}
