package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	maxOpenConns = 25
	maxIdleConns = 25
	maxLifeTime  = 100 * time.Second // max connection * seconds
)

// Client ... client.
type Client struct {
	DB *sql.DB
}

// NewClient ... get postgresql client.
func NewClient(
	host,
	user,
	password,
	database string,
	logEnable bool,
) *Client {
	dsnParts := []string{host}
	dsnParts = append(dsnParts,
		fmt.Sprintf("user=%s", user),
		fmt.Sprintf("password=%s", password),
		fmt.Sprintf("dbname=%s", database),
	)
	dbs := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		user,
		password,
		host,
		database)
	db, err := sql.Open("postgres", dbs)
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
	boil.DebugMode = logEnable
	return &Client{
		DB: db,
	}
}
