package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/lib/log"

	"github.com/go-chi/chi"
)

// GetURLParam ... get URL parameters from request
func GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// GetFormValue ... get form values from request
func GetFormValue(r *http.Request, key string) string {
	return r.FormValue(key)
}

// GetJSON ... get json data from request
func GetJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)
	if err != nil {
		ctx := r.Context()
		log.Errorf(ctx, "dec.Decode error: %s", err.Error())
		return err
	}
	return nil
}

// GetFormFile ... get file data from request
func GetFormFile(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	return r.FormFile(key)
}
