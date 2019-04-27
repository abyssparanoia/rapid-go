package mysql

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/src/lib/log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql" // import mysql
	"github.com/jmoiron/sqlx"
)

// NewSQLClient ... get mysql client
func NewSQLClient(cfg *SQLConfig) *sqlx.DB {
	ds := fmt.Sprintf("%s:%s@%s/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.ConnectionName,
		cfg.Database)

	db, err := sqlx.Open("mysql", ds)
	if err != nil {
		panic(err)
	}

	return db
}

// DumpSelectQuery ... output select query
func DumpSelectQuery(ctx context.Context, query sq.SelectBuilder) {
	qs, args, err := query.ToSql()
	if err != nil {
		log.Errorf(ctx, "DumpSelectQuery error: %s", err.Error())
		return
	}
	dumpQuery(ctx, qs, args)
}

// DumpInsertQuery ... output insert query
func DumpInsertQuery(ctx context.Context, query sq.InsertBuilder) {
	qs, args, err := query.ToSql()
	if err != nil {
		log.Errorf(ctx, "DumpInsertQuery error: %s", err.Error())
		return
	}
	dumpQuery(ctx, qs, args)
}

// DumpUpdateQuery ... output update query
func DumpUpdateQuery(ctx context.Context, query sq.UpdateBuilder) {
	qs, args, err := query.ToSql()
	if err != nil {
		log.Errorf(ctx, "DumpUpdateQuery error: %s", err.Error())
		return
	}
	dumpQuery(ctx, qs, args)
}

// DumpDeleteQuery ... output delete query
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
