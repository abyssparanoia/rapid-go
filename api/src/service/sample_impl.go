package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/lib/util"
	"github.com/abyssparanoia/rapid-go/api/src/model"
	"github.com/abyssparanoia/rapid-go/api/src/repository"
)

type sample struct {
	repo repository.Sample
}

func (s *sample) GetAll(ctx context.Context) (model.Sample, error) {
	log.Debugf(ctx, "call service beego")
	return model.Sample{
		ID:        123,
		Category:  "hoge",
		Name:      "sample太郎",
		Enabled:   true,
		CreatedAt: util.TimeNow(),
	}, nil
}

// NewSample ... サンプルサービスを取得する
func NewSample(repo repository.Sample) Sample {
	return &sample{
		repo: repo,
	}
}
