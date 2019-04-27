package httpheader

import (
	"context"
	"net/http"
)

// Service ... interface for http header
type Service interface {
	Get(ctx context.Context, r *http.Request) (Params, error)
}
