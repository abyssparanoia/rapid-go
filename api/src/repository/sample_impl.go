package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/api/src/lib/httpclient"
	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/lib/mysql"
	"github.com/abyssparanoia/rapid-go/api/src/lib/util"
	"github.com/abyssparanoia/rapid-go/api/src/model"

	sq "github.com/Masterminds/squirrel"
	_ "go.mercari.io/datastore/aedatastore" // mercari/datastoreの初期化
)

type sample struct {
	sql *sql.DB
}

func (r *sample) MySQLGet(ctx context.Context, id int64) (*model.Sample, error) {
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

	mysql.DumpSelectQuery(ctx, q)

	row := q.RunWith(r.sql).QueryRowContext(ctx)
	err := row.Scan(
		&ret.ID,
		&ret.Category,
		&ret.Name,
		&ret.Enabled,
		&ret.CreatedAt,
		&ret.UpdatedAt)
	if err != nil {
		log.Errorf(ctx, "MySQLGet: %s", err.Error())
		return ret, err
	}

	return ret, nil
}

func (r *sample) MySQLGetMulti(ctx context.Context, ids []int64) ([]*model.Sample, error) {
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

	mysql.DumpSelectQuery(ctx, q)

	rows, err := q.RunWith(r.sql).QueryContext(ctx)
	if err != nil {
		log.Errorf(ctx, "MySQLGetMulti: %s", err.Error())
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
			log.Errorf(ctx, "MySQLGet: %s", err.Error())
			rows.Close()
			return rets, err
		}
		rets = append(rets, ret)
	}

	return rets, nil
}

func (r *sample) MySQLInsert(ctx context.Context, obj *model.Sample) error {
	now := util.TimeNow()

	q := sq.Insert("sample").
		Columns("id", "category", "name", "enabled", "created_at", "updated_at").
		Values(obj.ID, obj.Category, obj.Name, 1, now, now)

	mysql.DumpInsertQuery(ctx, q)

	_, err := q.RunWith(r.sql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "MySQLInsert: %s", err.Error())
		return err
	}

	return nil
}

func (r *sample) MySQLUpdate(ctx context.Context, obj *model.Sample) error {
	now := util.TimeNow()

	q := sq.Update("sample").
		Set("name", obj.Name).
		Set("category", obj.Category).
		Set("enabled", obj.Enabled).
		Set("updated_at", now).
		Where(sq.Eq{"id": obj.ID})

	mysql.DumpUpdateQuery(ctx, q)

	res, err := q.RunWith(r.sql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "MySQLUpdate: %s", err.Error())
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		err = fmt.Errorf("no affected id = %d", obj.ID)
		log.Errorf(ctx, "MySQLUpdate: %s", err.Error())
		return err
	}

	return nil
}

func (r *sample) MySQLUpsert(ctx context.Context, obj *model.Sample) error {
	now := util.TimeNow()

	q := sq.Insert("sample").
		Columns("id", "name", "category", "enabled", "created_at", "updated_at").
		Values(obj.ID, obj.Category, obj.Name, 1, now, now).
		Suffix("ON DUPLICATE KEY UPDATE name = VALUES(name), updated_at = VALUES(updated_at)")

	mysql.DumpInsertQuery(ctx, q)

	_, err := q.RunWith(r.sql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "MySQLUpsert: %s", err.Error())
		return err
	}

	return nil
}

func (r *sample) MySQLDelete(ctx context.Context, id int64) error {
	q := sq.Delete("sample").Where(sq.Eq{"id": id})

	mysql.DumpDeleteQuery(ctx, q)

	res, err := q.RunWith(r.sql).ExecContext(ctx)
	if err != nil {
		log.Errorf(ctx, "MySQLDelete: %s", err.Error())
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		err = fmt.Errorf("no affected id = %d", id)
		log.Errorf(ctx, "MySQLDelete: %s", err.Error())
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
func NewSample(sql *sql.DB) Sample {
	return &sample{
		sql: sql,
	}
}
