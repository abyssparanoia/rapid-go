package renderer

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/errcode"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/unrolled/render"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// HandleError ... handle http error
func HandleError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	code, ok := errcode.Get(err)
	if !ok {
		Error(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	switch code {
	case http.StatusBadRequest:
		msg := fmt.Sprintf("%d StatusBadRequest: %s, %s", code, msg, err.Error())
		log.Warningf(ctx, msg)
	case http.StatusUnauthorized:
		msg := fmt.Sprintf("%d Unauthorized: %s, %s", code, msg, err.Error())
		log.Warningf(ctx, msg)
	case http.StatusForbidden:
		msg := fmt.Sprintf("%d Forbidden: %s, %s", code, msg, err.Error())
		log.Warningf(ctx, msg)
	case http.StatusNotFound:
		msg := fmt.Sprintf("%d NotFound: %s, %s", code, msg, err.Error())
		log.Warningf(ctx, msg)
	default:
		msg := fmt.Sprintf("%d: %s, %s", code, msg, err.Error())
		log.Errorf(ctx, msg)
	}

	Error(ctx, w, code, err.Error())
}

// Success ... render success response
func Success(ctx context.Context, w http.ResponseWriter) {
	r := render.New()
	r.JSON(w, http.StatusOK, NewResponseOK(http.StatusOK))
}

// Error ... render error response
func Error(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	r := render.New()
	r.JSON(w, status, NewResponseError(status, msg))
}

// JSON ... render json response
func JSON(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	r := render.New()
	r.JSON(w, status, v)
}

// HTML ... render html
func HTML(ctx context.Context, w http.ResponseWriter, status int, name string, values interface{}) {
	r := render.New()
	r.HTML(w, status, name, values)
}

// Text ... render text
func Text(ctx context.Context, w http.ResponseWriter, status int, body string) {
	r := render.New()
	r.Text(w, status, body)
}

// CSV ... render csv
func CSV(ctx context.Context, w http.ResponseWriter, name string, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", name))

	writer := csv.NewWriter(transform.NewWriter(w, japanese.ShiftJIS.NewEncoder()))
	for _, datum := range data {
		writer.Write(datum)
	}
	writer.Flush()
}
