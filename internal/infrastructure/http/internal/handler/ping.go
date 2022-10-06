package handler

import (
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pongpong"))
	if err != nil {
		panic(err)
	}
}
