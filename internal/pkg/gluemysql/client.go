package gluemysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/boil"
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
	boil.SetLocation(time.Local)
	boil.SetDB(db)
	return &Client{
		DB: db,
	}
}
