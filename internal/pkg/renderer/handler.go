package renderer

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/abyssparanoia/rapid-go/internal/pkg/error/httperror"
	"github.com/unrolled/render"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// HandleError ... handle http error
func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	statusCode := httperror.ErrToCode(err)
	ctxzap.AddFields(ctx, zap.Error(err))
	Error(ctx, w, statusCode, err.Error())
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
