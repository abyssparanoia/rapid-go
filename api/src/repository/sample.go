package repository

import (
	"context"

	"github.com/abyssparanoia/rapid-go/api/src/model"
)

// Sample ... リポジトリのインターフェース
type Sample interface {

	// MySQL
	MySQLGet(ctx context.Context, id int64) (*model.Sample, error)
	MySQLGetMulti(ctx context.Context, ids []int64) ([]*model.Sample, error)
	MySQLInsert(ctx context.Context, obj *model.Sample) error
	MySQLUpdate(ctx context.Context, obj *model.Sample) error
	MySQLUpsert(ctx context.Context, obj *model.Sample) error
	MySQLDelete(ctx context.Context, id int64) error

	// HTTP
	HTTPGet(ctx context.Context) error
	HTTPPost(ctx context.Context) error
}
