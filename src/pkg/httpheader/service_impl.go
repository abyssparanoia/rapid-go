package httpheader

import (
	"context"
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/pkg/log"

	validator "gopkg.in/go-playground/validator.v9"
)

const (
	headerKeySample string = "X-Sample"
)

type service struct {
}

func (s *service) Get(ctx context.Context, r *http.Request) (Params, error) {
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

// NewService ... get service
func NewService() Service {
	return &service{}
}
