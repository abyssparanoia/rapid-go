package basicauth

import (
	"net/http"
)

// Middleware ... ベーシック認証機能を提供するミドルウェア
type Middleware struct {
	Account *Account
}

// Handle ... ハンドラ
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

// NewMiddleware ... Middlewareを作成する
func NewMiddleware(account *Account) *Middleware {
	return &Middleware{
		Account: account,
	}
}
