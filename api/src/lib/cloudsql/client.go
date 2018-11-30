package cloudsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/abyssparanoia/gke-beego/api/src/lib/log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql" // MySQL Driverの読み込み
)

// NewCSQLClient ... CloudSQLのクライアントを取得する
func NewCSQLClient(cfg *CSQLConfig) *sql.DB {
	ds := fmt.Sprintf("%s:%s@cloudsql(%s)/",
		cfg.User,
		cfg.Password,
		cfg.ConnectionName)

	cli, err := sql.Open("mysql", ds)
	if err != nil {
		panic(err)
	}

	return cli
}

// DumpSelectQuery ... SELECTクエリを出力
func DumpSelectQuery(ctx context.Context, query sq.SelectBuilder) {
	qs, args, err := query.ToSql()
	if err != nil {
		log.Errorf(ctx, "DumpSelectQuery error: %s", err.Error())
		return
	}
	dumpQuery(ctx, qs, args)
}

// DumpInsertQuery ... INSERTクエリを出力
func DumpInsertQuery(ctx context.Context, query sq.InsertBuilder) {
	qs, args, err := query.ToSql()
	if err != nil {
		log.Errorf(ctx, "DumpInsertQuery error: %s", err.Error())
		return
	}
	dumpQuery(ctx, qs, args)
}

// DumpUpdateQuery ... UPDATEクエリを出力
func DumpUpdateQuery(ctx context.Context, query sq.UpdateBuilder) {
	qs, args, err := query.ToSql()
	if err != nil {
		log.Errorf(ctx, "DumpUpdateQuery error: %s", err.Error())
		return
	}
	dumpQuery(ctx, qs, args)
}

// DumpDeleteQuery ... DELETEクエリを出力
func DumpDeleteQuery(ctx context.Context, query sq.DeleteBuilder) {
	qs, args, err := query.ToSql()
	if err != nil {
		log.Errorf(ctx, "DumpDeleteQuery error: %s", err.Error())
		return
	}
	dumpQuery(ctx, qs, args)
}

func dumpQuery(ctx context.Context, queryString string, args []interface{}) {
	msg := fmt.Sprintf("[SQL Dump] %s, %s", queryString, args)
	log.Debugf(ctx, msg)
}
