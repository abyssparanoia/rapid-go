package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/api/src/model"
)

// Sample ... サービスのインターフェース
type Sample interface {
	GetAll(ctx context.Context) (model.Sample, error)
}
