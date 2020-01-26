package httpheader

import (
	"context"
	"net/http"
)

type dummy struct {
}

func (s *dummy) Get(ctx context.Context, r *http.Request) (Params, error) {
	h := Params{
		Sample: "sample",
	}
	return h, nil
}

// NewDummy ... get dummy
func NewDummy() Httpheader {
	return &dummy{}
}
