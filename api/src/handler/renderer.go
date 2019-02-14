package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/api/src/lib/log"

	"github.com/unrolled/render"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// HandleError ... 一番典型的なエラーハンドリング
func HandleError(ctx context.Context, w http.ResponseWriter, status int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Errorf(ctx, msg)
	RenderError(w, status, msg)
}

// RenderSuccess ... 成功レスポンスをレンダリングする
func RenderSuccess(w http.ResponseWriter) {
	r := render.New()
	r.Text(w, http.StatusOK, "success")
}

// RenderError ... エラーレスポンスをレンダリングする
func RenderError(w http.ResponseWriter, status int, msg string) {
	r := render.New()
	r.Text(w, status, fmt.Sprintf("%d %s", status, msg))
}

// RenderJSON ... JSONをレンダリングする
func RenderJSON(w http.ResponseWriter, status int, v interface{}) {
	r := render.New(render.Options{IndentJSON: true})
	r.JSON(w, status, v)
}

// RenderHTML ... HTMLをレンダリングする
func RenderHTML(w http.ResponseWriter, status int, name string, values interface{}) {
	r := render.New()
	r.HTML(w, status, name, values)
}

// RenderCSV ... CSVをレンダリングする
func RenderCSV(w http.ResponseWriter, name string, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", name))

	writer := csv.NewWriter(transform.NewWriter(w, japanese.ShiftJIS.NewEncoder()))
	for _, datum := range data {
		writer.Write(datum)
	}
	writer.Flush()
}
