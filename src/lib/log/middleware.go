package log

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/config"
)

// Middleware ... loggerを生成し、 リクエストログ掃き出し用のmiddleware
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		ctx := r.Context()
		NewLogger(ctx, config.IsEnvDeveloping())
		Infof(ctx, "Remote: %s [%s] %s", rAddr, method, path)
		next.ServeHTTP(w, r)
	})
}
