package httpheader

import (
	"context"
	"net/http"
)

// Httpheader ... interface for http header
type Httpheader interface {
	Get(ctx context.Context, r *http.Request) (Params, error)
}
