package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/gke-beego/api/src/lib/cloudsql"
	"github.com/abyssparanoia/gke-beego/api/src/lib/httpclient"
	"github.com/abyssparanoia/gke-beego/api/src/lib/log"
	"github.com/abyssparanoia/gke-beego/api/src/lib/util"
	"github.com/abyssparanoia/gke-beego/api/src/model"

	sq "github.com/Masterminds/squirrel"
	"go.mercari.io/datastore"
	_ "go.mercari.io/datastore/aedatastore" // mercari/datastoreの初期化
	"go.mercari.io/datastore/boom"
)

type sample struct {
	csql *sql.DB
}

// DataStore
func (r *sample) DataStoreGet(ctx context.Context, id int64) (*model.Sample, error) {
	dst := &model.Sample{
		ID: id,
	}
	b, err := boom.FromContext(ctx)
	if err != nil {
		log.Errorf(ctx, "boom from context error: %s", err.Error())
		return dst, err
	}
	if err := b.Get(dst); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, err
		}
		log.Errorf(ctx, "get error: %s", err.Error())
		return nil, err
	}
	return dst, nil
}

func (r *sample) DataStoreGetMulti(ctx context.Context, ids []int64) ([]*model.Sample, error) {
	ret := []*model.Sample{}
	b, err := boom.FromContext(ctx)
	if err != nil {
		log.Errorf(ctx, "boom from context error: %s", err.Error())
		return ret, err
	}
	bt := b.Batch()
	for _, id := range ids {
		dst := &model.Sample{
			ID: id,
		}
		bt.Get(dst, func(err error) error {
			if err != nil {
				if err == datastore.ErrNoSuchEntity {
					return nil
				}
				log.Errorf(ctx, "get multi error: %s, id: %d", err.Error(), dst.ID)
				return err
			}
			ret = append(ret, dst)
			return nil
		})
	}
	err = bt.Exec()
	if err != nil {
		log.Errorf(ctx, "get multi error: %s", err.Error())
		return nil, err
	}
	return ret, nil
}

func (r *sample) DataStoreGetByQuery(ctx context.Context, category string) ([]*model.Sample, error) {
	ret := []*model.Sample{}
	b, err := boom.FromContext(ctx)
	if err != nil {
		log.Errorf(ctx, "boom from context error: %s", err.Error())
		return ret, err
	}
	q := b.NewQuery("Sample").Filter("Category =", category).Filter("Enabled =", true).Order("-CreatedAt")
	_, err = b.GetAll(q, &ret)
	if err != nil {
		log.Errorf(ctx, "get by query error: "+err.Error())
		return ret, err
	}
	return ret, nil
}

func (r *sample) DataStoreInsert(ctx context.Context, obj *model.Sample) (int64, error) {

	return 0, nil
}

func (r *sample) DataStoreInsertMulti(ctx context.Context, objs []*model.Sample) ([]int64, error) {
	return []int64{}, nil
}

func (r *sample) DataStoreUpdate(ctx context.Context, obj *model.Sample) (int64, error) {
	return 0, nil
}

func (r *sample) DataStoreUpdateMulti(ctx context.Context, objs []*model.Sample) ([]int64, error) {
	return []int64{}, nil
}

func (r *sample) DataStoreUpsert(ctx context.Context, src *model.Sample) (int64, error) {
	var id int64
	b, err := boom.FromContext(ctx)
	if err != nil {
		log.Errorf(ctx, "boom from context error: %s", err.Error())
		return id, err
	}
	key, err := b.Put(src)
	if err != nil {
		log.Errorf(ctx, "upsert error: %s", err.Error())
		return id, err
	}
	id = key.ID()
	return id, nil
}

func (r *sample) DataStoreUpsertMulti(ctx context.Context, srcs []*model.Sample) ([]int64, error) {
	ids := []int64{}
	b, err := boom.FromContext(ctx)
	if err != nil {
		log.Errorf(ctx, "boom from context error: %s", err.Error())
		return ids, err
	}
	bt := b.Batch()
	for _, src := range srcs {
		bt.Put(src, func(key datastore.Key, err error) error {
			if err != nil {
				log.Errorf(ctx, "upsert error: %s, id: %d", err.Error(), key.ID())
				return err
			}
			ids = append(ids, key.ID())
			return nil
		})
	}
	err = bt.Exec()
	if err != nil {
		log.Errorf(ctx, "upsert error: %s", err.Error())
		return ids, err
	}
	return ids, nil
}

func (r *sample) DataStoreDelete(ctx context.Context, id int64) (int64, error) {
	return 0, nil
}

func (r *sample) DataStoreDeleteMulti(ctx context.Context, id int64) ([]int64, error) {
	return nil, nil
}

// CloudSQL
func (r *sample) CloudSQLGet(ctx context.Context, id int64) (*model.Sample, error) {
	var ret *model.Sample

	q := sq.Select(
		"id",
		"category",
		"name",
		"enabled",
		"created_at",
		"updated_at").
		From("sample").
		Where(sq.Eq{
			"id":      id,
			"enabled": 1,
		})

	cloudsql.DumpSelectQuery(ctx, q)

	row := q.RunWith(r.csql).QueryRowContext(ctx)
	err := row.Scan(
		&ret.ID,
		&ret.Category,
		&ret.Name,
		&ret.Enabled,
		&ret.CreatedAt,
		&ret.UpdatedAt)
	if err != nil {
		log.Errorf(ctx, "CloudSQLGet: %s", err.Error())
		return ret, err
	}

	return ret, nil
}

func (r *sample) CloudSQLGetMulti(ctx context.Context, ids []int64) ([]*model.Sample, error) {
	var rets []*model.Sample

	q := sq.Select(
		"id",
		"name",
		"category",
		"enabled",
		"created_at",
		"updated_at").
		From("sample").
		Where(sq.Eq{
			"id":      ids,
			"enabled": 1,
		})

	cloudsql.DumpSelectQuery(ctx, q)

	rows, err := q.RunWith(r.csql).QueryContext(ctx)
	if err != nil {
		log.Errorf(ctx, "CloudSQLGetMulti: %s", err.Error())
		return rets, err
	}

	for rows.Next() {
		var ret *model.Sample
		err := rows.Scan(
			&ret.ID,
			&ret.Name,
			&ret.Category,
			&ret.Enabled,
			&ret.CreatedAt,
			&ret.UpdatedAt)
		if err != nil {
			log.Errorf(ctx, "CloudSQLGet: %s", err.Error())
			rows.Close()
			return rets, err
		}
		rets = append(rets, ret)
	}

	return rets, nil
}

func (r *sample) CloudSQLInsert(ctx context.Context, obj *model.Sample) error {
	now := util.TimeNow()

	q := sq.Insert("sample").
		Columns("id", "category", "name", "enabled", "created_at", "updated_at").
		Values(obj.ID, obj.Category, obj.Name, 1, now, now)

	cloudsql.DumpInsertQuery(ctx, q)

	_, err := q.RunWith(r.csql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "CloudSQLInsert: %s", err.Error())
		return err
	}

	return nil
}

func (r *sample) CloudSQLUpdate(ctx context.Context, obj *model.Sample) error {
	now := util.TimeNow()

	q := sq.Update("sample").
		Set("name", obj.Name).
		Set("category", obj.Category).
		Set("enabled", obj.Enabled).
		Set("updated_at", now).
		Where(sq.Eq{"id": obj.ID})

	cloudsql.DumpUpdateQuery(ctx, q)

	res, err := q.RunWith(r.csql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "CloudSQLUpdate: %s", err.Error())
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		err = fmt.Errorf("no affected id = %d", obj.ID)
		log.Errorf(ctx, "CloudSQLUpdate: %s", err.Error())
		return err
	}

	return nil
}

func (r *sample) CloudSQLUpsert(ctx context.Context, obj *model.Sample) error {
	now := util.TimeNow()

	q := sq.Insert("sample").
		Columns("id", "name", "category", "enabled", "created_at", "updated_at").
		Values(obj.ID, obj.Category, obj.Name, 1, now, now).
		Suffix("ON DUPLICATE KEY UPDATE name = VALUES(name), updated_at = VALUES(updated_at)")

	cloudsql.DumpInsertQuery(ctx, q)

	_, err := q.RunWith(r.csql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "CloudSQLUpsert: %s", err.Error())
		return err
	}

	return nil
}

func (r *sample) CloudSQLDelete(ctx context.Context, id int64) error {
	q := sq.Delete("sample").Where(sq.Eq{"id": id})

	cloudsql.DumpDeleteQuery(ctx, q)

	res, err := q.RunWith(r.csql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "CloudSQLDelete: %s", err.Error())
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		err = fmt.Errorf("no affected id = %d", id)
		log.Errorf(ctx, "CloudSQLDelete: %s", err.Error())
		return err
	}

	return nil
}

// HTTP
func (r *sample) HTTPGet(ctx context.Context) error {
	status, body, err := httpclient.Get(ctx, "https://www.google.co.jp/", nil)
	if err != nil {
		log.Errorf(ctx, "HTTPGet: %s", err.Error())
		return err
	}
	if status != http.StatusOK {
		err := fmt.Errorf("http status: %d", status)
		return err
	}
	str := util.BytesToStr(body)
	log.Debugf(ctx, "body length: %d", len(str))
	return nil
}

func (r *sample) HTTPPost(ctx context.Context) error {
	return nil
}

// NewSample ... サンプルリポジトリを取得する
func NewSample(csql *sql.DB) Sample {
	return &sample{
		csql: csql,
	}
}
