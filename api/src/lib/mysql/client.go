package mysql

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/jinzhu/gorm"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql" // MySQL Driverの読み込み
)

// NewSQLClient ... MySQLのクライアントを取得する
func NewSQLClient(cfg *SQLConfig) *gorm.DB {
	ds := fmt.Sprintf("%s:%s@%s/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.ConnectionName,
		cfg.Database)

	db, err := gorm.Open("mysql", ds)
	if err != nil {
		panic(err)
	}

	return db
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
