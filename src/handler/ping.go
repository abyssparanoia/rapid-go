package handler

import (
	"net/http"
)

// Ping ... 生存確認
func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pongpong"))
}
