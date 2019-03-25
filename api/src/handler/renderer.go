package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/api/src/lib/errcode"
	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/model"

	"github.com/unrolled/render"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// HandleError ... 一番典型的なエラーハンドリング
func HandleError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	code, ok := errcode.Get(err)
	if !ok {
		RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	switch code {
	case http.StatusBadRequest:
		msg := fmt.Sprintf("%d StatusBadRequest: %s, %s", code, msg, err.Error())
		log.Warningf(ctx, msg)
		RenderError(w, code, msg)
	case http.StatusForbidden:
		msg := fmt.Sprintf("%d Forbidden: %s, %s", code, msg, err.Error())
		log.Warningf(ctx, msg)
		RenderError(w, code, msg)
	case http.StatusNotFound:
		msg := fmt.Sprintf("%d NotFound: %s, %s", code, msg, err.Error())
		log.Warningf(ctx, msg)
		RenderError(w, code, msg)
	default:
		msg := fmt.Sprintf("%d: %s, %s", code, msg, err.Error())
		log.Errorf(ctx, msg)
		RenderError(w, code, msg)
	}
}

// RenderSuccess ... 成功レスポンスをレンダリングする
func RenderSuccess(w http.ResponseWriter) {
	r := render.New()
	r.JSON(w, http.StatusOK, model.NewResponseOK(http.StatusOK))
}

// RenderError ... エラーレスポンスをレンダリングする
func RenderError(w http.ResponseWriter, status int, msg string) {
	r := render.New()
	r.JSON(w, status, model.NewResponseError(status, msg))
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
