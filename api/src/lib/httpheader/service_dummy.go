package httpheader

import (
	"context"
	"net/http"
)

type dummyService struct {
}

func (s *dummyService) Get(ctx context.Context, r *http.Request) (Params, error) {
	h := Params{
		// EDIT: ここに任意のダミーヘッダーを入れる
		Sample: "sample",
	}
	return h, nil
}

// NewDummyService ... DummyServiceを作成する
func NewDummyService() Service {
	return &dummyService{}
}
