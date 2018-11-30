package httpheader

import (
	"context"
	"net/http"
)

// Service ... HTTPHeaderに関する機能を提供する
type Service interface {
	Get(ctx context.Context, r *http.Request) (Params, error)
}
