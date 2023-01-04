package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	maxOpenConns = 25
	maxIdleConns = 25
	maxLifeTime  = 100 * time.Second // max conection * seconds
)

// Client ... client
type Client struct {
	DB *sql.DB
}

// NewClient ... get mysql client
func NewClient(
	host,
	user,
	password,
	database string,
) *Client {
	dbs := fmt.Sprintf("%s:%s@%s/%s?parseTime=true&charset=utf8mb4",
		user,
		password,
		host,
		database)
	db, err := sql.Open("mysql", dbs)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxLifeTime)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	boil.SetLocation(time.Local)
	boil.SetDB(db)
	boil.DebugMode = true
	return &Client{
		DB: db,
	}
}
