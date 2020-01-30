package handler

import (
	"net/http"
)

// Ping ... confirmation survival
func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pongpong"))
}
