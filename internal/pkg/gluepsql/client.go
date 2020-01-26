package gluepsql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewClient ... get psql client
func NewClient(cfg *Config) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}
