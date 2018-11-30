package log

import (
	"log"
	"net/http"
)

// Log ... リクエストログ掃き出し用のmiddleware
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote: %s [%s] %s", rAddr, method, path)
		h.ServeHTTP(w, r)
	})
}
