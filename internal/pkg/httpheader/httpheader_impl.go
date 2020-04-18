package httpheader

import (
	"context"
	"net/http"

	validator "github.com/go-playground/validator"
)

const (
	headerKeySample string = "X-Sample"
)

type httpheader struct {
}

func (s *httpheader) Get(ctx context.Context, r *http.Request) (Params, error) {
	h := Params{
		Sample: r.Header.Get(headerKeySample),
	}

	v := validator.New()
	if err := v.Struct(h); err != nil {
		return h, err
	}

	return h, nil
}

// New ... get httpheader
func New() Httpheader {
	return &httpheader{}
}
