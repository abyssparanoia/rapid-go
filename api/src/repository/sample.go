package repository

import (
	"context"

	"github.com/abyssparanoia/gke-beego/api/src/model"
)

// Sample ... リポジトリのインターフェース
type Sample interface {
	// DataStore
	DataStoreGet(ctx context.Context, id int64) (*model.Sample, error)
	DataStoreGetMulti(ctx context.Context, ids []int64) ([]*model.Sample, error)
	DataStoreGetByQuery(ctx context.Context, category string) ([]*model.Sample, error)
	DataStoreInsert(ctx context.Context, obj *model.Sample) (int64, error)
	DataStoreInsertMulti(ctx context.Context, objs []*model.Sample) ([]int64, error)
	DataStoreUpdate(ctx context.Context, obj *model.Sample) (int64, error)
	DataStoreUpdateMulti(ctx context.Context, objs []*model.Sample) ([]int64, error)
	DataStoreUpsert(ctx context.Context, obj *model.Sample) (int64, error)
	DataStoreUpsertMulti(ctx context.Context, objs []*model.Sample) ([]int64, error)
	DataStoreDelete(ctx context.Context, id int64) (int64, error)
	DataStoreDeleteMulti(ctx context.Context, id int64) ([]int64, error)

	// CloudSQL
	CloudSQLGet(ctx context.Context, id int64) (*model.Sample, error)
	CloudSQLGetMulti(ctx context.Context, ids []int64) ([]*model.Sample, error)
	CloudSQLInsert(ctx context.Context, obj *model.Sample) error
	CloudSQLUpdate(ctx context.Context, obj *model.Sample) error
	CloudSQLUpsert(ctx context.Context, obj *model.Sample) error
	CloudSQLDelete(ctx context.Context, id int64) error

	// HTTP
	HTTPGet(ctx context.Context) error
	HTTPPost(ctx context.Context) error
}
