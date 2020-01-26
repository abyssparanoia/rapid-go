package httpheader

import (
	"context"
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"

	validator "gopkg.in/go-playground/validator.v9"
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
		log.Warningf(ctx, "Header param validation error: %s", err.Error())
		return h, err
	}

	return h, nil
}

// New ... get httpheader
func New() Httpheader {
	return &httpheader{}
}
