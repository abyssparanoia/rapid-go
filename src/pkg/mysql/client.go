package mysql

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Client ... client
type Client struct {
	db *gorm.DB
}

// GetDB ... get db with set logger
func (c *Client) GetDB(ctx context.Context) *gorm.DB {
	db := c.db.New()
	db.SetLogger(gorm.Logger{
		LogWriter: NewLogger(ctx),
	})
	return db
}

// NewClient ... get gorm client
func NewClient(cfg *Config) *Client {
	dbs := fmt.Sprintf("%s:%s@%s/%s?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DB)
	db, err := gorm.Open("mysql", dbs)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	return &Client{
		db: db,
	}
}
