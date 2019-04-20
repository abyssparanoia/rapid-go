package log

import (
	"net/http"
)

// Middleware ... output request log
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		ctx := r.Context()
		ctx = NewLogger(ctx)
		Infof(ctx, "Remote: %s [%s] %s", rAddr, method, path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
