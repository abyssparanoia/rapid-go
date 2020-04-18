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
	db *sql.DB
}

// NewClient ... get gorm client
func NewClient(cfg *Config) *Client {
	dbs := fmt.Sprintf("%s:%s@%s/%s?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DB)
	db, err := sql.Open("mysql", dbs)
	if err != nil {
		panic(err)
	}
	boil.SetLocation(time.Local)
	boil.SetDB(db)
	return &Client{
		db: db,
	}
}
