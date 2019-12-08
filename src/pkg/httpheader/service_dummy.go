package httpheader

import (
	"context"
	"net/http"
)

type dummyService struct {
}

func (s *dummyService) Get(ctx context.Context, r *http.Request) (Params, error) {
	h := Params{
		Sample: "sample",
	}
	return h, nil
}

// NewDummyService ... get dummy service
func NewDummyService() Service {
	return &dummyService{}
}
